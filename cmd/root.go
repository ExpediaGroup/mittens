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
	"log"
	"mittens/cmd/flags"
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
		log.Printf("invalid HTTP options: %v", err)
		validationError = true
	}
	grpcRequests, err := opts.GetWarmupGrpcRequests()
	if err != nil {
		log.Printf("invalid grpc options: %v", err)
		validationError = true
	}
	targetOptions, err := opts.GetWarmupTargetOptions()
	if err != nil {
		log.Printf("invalid target options: %v", err)
		validationError = true
	}

	httpHeaders := opts.GetWarmupHTTPHeaders()

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

				if opts.FetchAuthToken {
					token, err := target.FetchAuthToken()
					if err != nil {
						log.Printf("FetchAuthToken call failed: %s", err.Error())
					} else {
						log.Printf("AuthToken: %s", token)
						AuthToken := "Authorization:" + token
						httpHeaders = append(httpHeaders, AuthToken)
					}
				}

				elapsed := time.Since(start).Seconds()

				log.Printf("💚 Target took %d second(s) to become ready", int(elapsed))

				globalMaxDurationSecondsLeft := opts.MaxDurationSeconds - int(elapsed)

				maxDurationInSeconds := Min(globalMaxDurationSecondsLeft, opts.MaxWarmupDurationSeconds)

				if maxDurationInSeconds < opts.MaxWarmupDurationSeconds {
					log.Printf("⚠️ Warmup requests will only run for %d seconds instead of the configured %d seconds as to meet the global maximum duration of %d seconds", maxDurationInSeconds, opts.MaxWarmupDurationSeconds, opts.MaxDurationSeconds)
				}

				wp := warmup.Warmup{
					Target:                   target,
					Concurrency:              opts.GetConcurrency(),
					HttpRequests:             httpRequests,
					GrpcRequests:             grpcRequests,
					HttpHeaders:              httpHeaders,
					RequestDelayMilliseconds: opts.RequestDelayMilliseconds,
					ConcurrencyTargetSeconds: opts.GetConcurrencyTargetSeconds(),
				}

				wp.Run(hasHttpRequests, hasGrpcRequests, maxDurationInSeconds, &requestsSentCounter)
			} else {
				log.Print("Target still not ready. Giving up!")
			}
		}
		c1 <- true
	})

	<-c1
	log.Println("🟢 Warmup completed")
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
		log.Print("🛑 Warmup did not run. Mittens readiness probe will fail 🙁")
	} else {
		if requestsSentCounter == 0 {
			log.Print("🛑 Warm up finished but no requests were sent 🙁")
		} else {
			log.Printf("Warm up finished 😊 Approximately %d reqs were sent", requestsSentCounter)
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
