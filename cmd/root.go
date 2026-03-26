//Copyright 2019 Expedia, Inc.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package cmd

import (
	"flag"
	"log/slog"
	"mittens/cmd/flags"
	"mittens/internal/pkg/logging"
	"mittens/internal/pkg/probe"
	"mittens/internal/pkg/safe"
	"mittens/internal/pkg/warmup"
	"os"
	"time"
)

var opts *flags.Root

// CreateConfig creates a flag set and parses the command line arguments.
func CreateConfig() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	opts = &flags.Root{}
	opts.InitFlags()
	flag.Parse()

	if err := logging.Setup(opts.LogLevel, opts.LogFormat, opts.LogDateTimeFormat); err != nil {
		slog.Error("Failed to configure logging, using defaults", "error", err)
	}
}

// RunCmdRoot runs the main logic
//
//	It blocks forever unless `-exit-after-warmup` is set to true
func RunCmdRoot() {
	requestsSent := safe.DoAndReturn(run, 0)
	postProcess(requestsSent)
	block()
}

// run runs the main logic and returns the number of warmup requests actually sent.
func run() int {
	if opts.FileProbe.Enabled {
		probe.WriteFile(opts.FileProbe.LivenessPath)
	}

	var validationError bool
	httpRequests, err := opts.GetWarmupHTTPRequests()
	if err != nil {
		slog.Error("Invalid HTTP options", "error", err)
		validationError = true
	}
	grpcRequests, err := opts.GetWarmupGrpcRequests()
	if err != nil {
		slog.Error("Invalid gRPC options", "error", err)
		validationError = true
	}
	targetOptions, err := opts.GetWarmupTargetOptions()
	if err != nil {
		slog.Error("Invalid target options", "error", err)
		validationError = true
	}

	// this is used to decide on whether we should create goroutines for HTTP and/or gRPC requests
	// since requests are passed to a channel after that point we need to store that info and pass it
	var hasHttpRequests bool
	if len(opts.HTTP.Requests) > 0 {
		hasHttpRequests = true
	}
	var hasGrpcRequests bool
	if len(opts.Grpc.Requests) > 0 {
		hasGrpcRequests = true
	}

	// The next block contains the "wait for target readiness" + "warmup" logic.
	c1 := make(chan bool, 1)

	requestsSentCounter := 0

	// current time
	start := time.Now()

	go safe.Do(func() {
		if !validationError {
			target := createTarget(targetOptions)

			maxReadinessWaitDurationInSeconds := Min(opts.MaxDurationSeconds, opts.MaxReadinessWaitSeconds)

			if err := target.WaitForReadinessProbe(maxReadinessWaitDurationInSeconds, opts.GetWarmupHTTPHeaders()); err == nil {
				elapsed := time.Since(start).Seconds()

				slog.Info("Target is ready", "elapsed_seconds", int(elapsed))

				globalMaxDurationSecondsLeft := opts.MaxDurationSeconds - int(elapsed)

				maxDurationInSeconds := Min(globalMaxDurationSecondsLeft, opts.MaxWarmupDurationSeconds)

				if maxDurationInSeconds < opts.MaxWarmupDurationSeconds {
					slog.Warn("Warmup duration capped by global maximum", "actual_seconds", maxDurationInSeconds, "configured_seconds", opts.MaxWarmupDurationSeconds, "max_duration_seconds", opts.MaxDurationSeconds)
				}

				wp := warmup.Warmup{
					Target:                   target,
					Concurrency:              opts.GetConcurrency(),
					HttpRequests:             httpRequests,
					GrpcRequests:             grpcRequests,
					HttpHeaders:              opts.GetWarmupHTTPHeaders(),
					RequestDelayMilliseconds: opts.RequestDelayMilliseconds,
					ConcurrencyTargetSeconds: opts.GetConcurrencyTargetSeconds(),
				}

				wp.Run(hasHttpRequests, hasGrpcRequests, maxDurationInSeconds, &requestsSentCounter)
			} else {
				slog.Error("Target still not ready. Giving up!")
			}
		}
		c1 <- true
	})

	<-c1
	slog.Info("Warmup completed")
	return requestsSentCounter
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// block blocks forever unless `-exit-after-warmup` is set to true
func block() {
	if !opts.ExitAfterWarmup {
		select {}
	}
}

// postProcess includes steps that run once the warmup finishes.
// For now this either announces that the app is ready or fails the readiness probe.
// The latter only happens if mittens did not send any requests and the user allows the readiness to fail.
func postProcess(requestsSentCounter int) {
	if opts.FailReadiness && requestsSentCounter == 0 {
		slog.Error("Warmup did not run. Mittens readiness probe will fail")
	} else {
		if requestsSentCounter == 0 {
			slog.Warn("Warmup finished but no requests were sent")
		} else {
			slog.Info("Warmup finished", "requests_sent", requestsSentCounter)
		}

		if opts.FileProbe.Enabled {
			probe.WriteFile(opts.FileProbe.ReadinessPath)
		}
	}
}

// createTarget creates the target versus which mittens will run.
func createTarget(targetOptions warmup.TargetOptions) warmup.Target {
	return warmup.NewTarget(
		opts.GetReadinessHTTPClient(),
		opts.GetReadinessGrpcClient(),
		opts.GetHTTPClient(),
		opts.GetGrpcClient(),
		targetOptions,
	)
}
