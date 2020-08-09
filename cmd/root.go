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
	"fmt"
	"log"
	"mittens/cmd/flags"
	"mittens/pkg/probe"
	"mittens/pkg/warmup"
	"os"
	"os/signal"
	"syscall"
)

var opts *flags.Root

// CreateConfig creates a flag set and parses the command line arguments.
func CreateConfig() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	opts = &flags.Root{}
	opts.InitFlags()
	flag.Parse()
}

// RunCmdRoot runs the main logic.
func RunCmdRoot() {
	var probeServer *probe.Server

	if opts.ServerProbe.Enabled {
		probeServer = startServerProbe(
			opts.ServerProbe.Port,
			opts.ServerProbe.LivenessPath,
			opts.ServerProbe.ReadinessPath,
		)
	}

	if opts.FileProbe.Enabled {
		probe.WriteFile(opts.FileProbe.LivenessPath)
	}

	var validationError bool
	httpRequests, err := opts.GetWarmupHTTPRequests()
	if err != nil {
		fmt.Errorf("invalid HTTP options: %v", err)
		validationError = true
	}
	grpcRequests, err := opts.GetWarmupGrpcRequests()
	if err != nil {
		fmt.Errorf("invalid grpc options: %v", err)
		validationError = true
	}
	targetOptions, err := opts.GetWarmupTargetOptions()
	if err != nil {
		fmt.Errorf("invalid target options: %v", err)
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

	if !validationError {
		requestsSentCounter := 0
		target := createTarget(targetOptions)
		if err := target.WaitForReadinessProbe(); err == nil {
			log.Print("💚 Target is ready")

			wp := warmup.Warmup{
				Target:                   target,
				MaxDurationSeconds:       opts.GetMaxDurationSeconds(),
				Concurrency:              opts.GetConcurrency(),
				HttpRequests:             httpRequests,
				GrpcRequests:             grpcRequests,
				HttpHeaders:              opts.GetWarmupHTTPHeaders(),
				GrpcHeaders:              opts.GetWarmupGrpcHeaders(),
				RequestDelayMilliseconds: opts.RequestDelayMilliseconds,
			}
			wp.Run(hasHttpRequests, hasGrpcRequests, &requestsSentCounter)
		} else {
			log.Print("Target still not ready. Giving up!")
		}

		postProcess(requestsSentCounter, probeServer)
	}

	// Block forever if we don't want to wait after the warmup finishes
	if !opts.ExitAfterWarmup {
		select {}
	}
}

// postProcess includes steps that run once the warmup finishes.
// For now this either announces that the app is ready or fails the readiness probe.
// The latter only happens if mittens did not send any requests and the user allows the readiness to fail.
func postProcess(requestsSentCounter int, probeServer *probe.Server) {
	if opts.FailReadiness && requestsSentCounter == 0 {
		log.Print("🛑 Warmup did not run. Mittens readiness probe will fail 🙁")
	} else {
		if requestsSentCounter == 0 {
			log.Print("🛑 Warm up finished but no requests were sent 🙁")
		} else {
			log.Printf("Warm up finished 😊 Approximately %d reqs were sent", requestsSentCounter)
		}

		if opts.ServerProbe.Enabled {
			probeServer.IsReady(true)
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

// startServerProbe starts a web server that can be used for readiness and liveness checks.
func startServerProbe(port int, livenessPath, readinessPath string) *probe.Server {

	serverErr := make(chan struct{})
	probeServer := probe.NewServer(port, livenessPath, readinessPath)
	go func() {
		if err := probeServer.ListenAndServe(); err != nil {
			if err.Error() != "HTTP: Server closed" {
				log.Printf("Probe server: %v", err)
				close(serverErr)
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-serverErr:
			log.Print("Received probe server error")
		case sig := <-sigs:
			log.Printf("Received %s signal", sig)
			probeServer.Shutdown()
		}
	}()
	return probeServer
}
