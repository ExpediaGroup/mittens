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
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"mittens/cmd/flags"
	"mittens/pkg/probe"
	"mittens/pkg/response"
	"mittens/pkg/warmup"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var opts *flags.Root

func CreateConfig() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var cfgFile string

	flag.StringVar(&cfgFile, "config", "", "Config file to be used. If empty configs will be read from cmd.")
	opts = &flags.Root{}
	opts.InitFlags()
	flag.Parse()

	if cfgFile != "" {
		log.Printf("Reading configs from file: %v", cfgFile)
		file, err := os.Open(cfgFile)
		if err != nil {
			log.Print("Can't open config file: ", err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&opts)
		if err != nil {
			log.Print("Can't decode config JSON: ", err)
		}
	}
}

func RunCmdRoot() {
	stop := make(chan struct{})
	done := make(chan struct{})
	var probeServer *probe.Server

	if opts.ServerProbe.Enabled {
		probeServer = start(
			opts.ServerProbe.Port,
			opts.ServerProbe.LivenessPath,
			opts.ServerProbe.ReadinessPath,
			stop,
			done,
		)
	}

	if opts.FileProbe.Enabled {
		probe.WriteFile(opts.FileProbe.LivenessPath)
	}

	targetOptions, err := opts.GetWarmupTargetOptions()

	var result warmup.Result
	if err == nil {
		wp, err1 := createWarmup(targetOptions, done)
		if err1 == nil {
			result = runWarmup(wp, done)
			log.Print("Warm up finished 😊")
		} else {
			log.Print("Target still not ready. Giving up! No requests were sent 🙁")
		}

		if opts.FailReadiness && err1 != nil {
			log.Print("Mittens is not ready. Warmup did not run. 🛑")
		} else {
			if opts.ServerProbe.Enabled {
				probeServer.IsReady(true)
			}
			if opts.FileProbe.Enabled {
				probe.WriteFile(opts.FileProbe.ReadinessPath)
			}
		}

	}
	log.Printf("Summary: SuccesfulRequests: %d, FailedRequests: %d", result.SuccesfulRequests, result.FailedRequests)

	if opts.ExitAfterWarmup {
		// exit after warmup, we close the stop/done channels
		// in case probe server is used the done channel is closed by the server to ensure graceful termination
		if opts.ServerProbe.Enabled {
			close(stop)
		} else {
			close(done)
		}
	} else {
		select {}
	}
	<-done
}

func runWarmup(wp warmup.Warmup, done chan struct{}) warmup.Result {
	rand.Seed(time.Now().UnixNano()) // initialize seed only once to prevent deterministic/repeated calls every time we run
	succesfulRequests := 0
	failedRequests := 0

	httpHeaders := opts.GetWarmupHttpHeaders()
	httpRequests, err := opts.GetWarmupHttpRequests(done)
	if err != nil {
		log.Printf("Http options: %v", err)
	}
	grpcHeaders := opts.GetWarmupGrpcHeaders()
	grpcRequests, err := opts.GetWarmupGrpcRequests(done)
	if err != nil {
		log.Printf("Grpc options: %v", err)
	}
	httpResponse := wp.HttpWarmup(httpHeaders, httpRequests)
	grpcResponse := wp.GrpcWarmup(grpcHeaders, grpcRequests)
	response := merge(httpResponse, grpcResponse)

	for r := range response {
		if r.Err != nil {
			log.Printf("🔴 %s response %d milliseconds: error: %v", r.Type, r.Duration/time.Millisecond, r.Err)
			failedRequests++
		} else {
			succesfulRequests++
			log.Printf("🟢 %s response %d milliseconds: OK", r.Type, r.Duration/time.Millisecond)
		}
	}
	return warmup.Result{SuccesfulRequests: succesfulRequests, FailedRequests: failedRequests}
}

// TODO: this is doing too many things including waiting for target to become ready. split into smaller blocks.
func createWarmup(targetOptions warmup.TargetOptions, done chan struct{}) (warmup.Warmup, error) {
	wp, err := warmup.NewWarmup(
		opts.GetReadinessHttpClient(),
		opts.GetReadinessGrpcClient(),
		opts.GetHttpClient(),
		opts.GetGrpcClient(),
		opts.GetWarmupOptions(),
		targetOptions,
		done,
	)
	return wp, err
}

func start(port int, livenessPath, readinessPath string, stop <-chan struct{}, done chan struct{}) *probe.Server {

	serverErr := make(chan struct{})
	probeServer := probe.NewServer(port, livenessPath, readinessPath)
	go func() {
		if err := probeServer.ListenAndServe(); err != nil {
			if err.Error() != "Http: Server closed" {
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
			close(done)
		case <-stop:
			log.Print("Received stop signal")
			probeServer.Shutdown()
			close(done)
		case sig := <-sigs:
			log.Printf("Received %s signal", sig)
			probeServer.Shutdown()
			close(done)
		}
	}()
	return probeServer
}

// 'fan in' see: https://blog.golang.org/pipelines
func merge(cs ...<-chan response.Response) <-chan response.Response {

	var wg sync.WaitGroup
	out := make(chan response.Response)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan response.Response) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
