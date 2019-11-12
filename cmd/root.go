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
	"mittens/cmd/flags"
	"mittens/pkg/probe"
	"mittens/pkg/warmup"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

var opts *flags.Root

func init() {
	var cfgFile string

	flag.StringVar(&cfgFile,"config", "","Config file to be used. If empty configs will be read from cmd.")
	opts = flags.NewRoot()
	opts.InitFlags()
	flag.Parse()

	if cfgFile != "" {
		log.Printf("reading configs from file: %v", cfgFile)
		file, err := os.Open(cfgFile)
		if err != nil {
			log.Fatal("can't open config file: ", err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&opts)
		if err != nil {
			log.Fatal("can't decode config JSON: ", err)
		}
	}
	flag.Parse()
}

func RunCmdRoot() {

	// CPU profile
	if opts.Profile.CPU != "" {

		log.Printf("CPU profile will be written to %s file", opts.Profile.CPU)
		// CPU profile
		f, err := os.Create(opts.Profile.CPU)
		if err != nil {
			log.Fatalf("could not create CPU profile: %v", err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("could not start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

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

	wp, err := warmup.NewWarmup(
		opts.GetHttpClient(),
		opts.GetGrpcClient(),
		opts.GetWarmupOptions(),
		opts.GetWarmupTargetOptions(),
		done,
	)

	if err != nil {
		log.Fatalf("new warmup: %v", err)
	}

	httpOptions := opts.GetWarmupHttpHeaders()
	httpRequests, err := opts.GetWarmupHttpRequests(done)
	if err != nil {
		log.Fatalf("http options: %v", err)
	}

	grpcOptions := opts.GetWarmupGrpcHeaders()
	grpcRequests, err := opts.GetWarmupGrpcRequests(done)
	if err != nil {
		log.Fatalf("grpc options: %v", err)
	}

	httpResponse := wp.HttpWarmup(httpOptions, httpRequests)
	grpcResponse := wp.GrpcWarmup(grpcOptions, grpcRequests)

	response := merge(httpResponse, grpcResponse)
	for r := range response {
		if r.Error != nil {
			log.Printf("%s response %d milliseconds: error: %v", r.Type, r.Duration/time.Millisecond, r.Error)
		}
		log.Printf("%s response %d milliseconds: OK", r.Type, r.Duration/time.Millisecond)
	}

	if opts.ServerProbe.Enabled {
		probeServer.IsReady(true)
	}
	if opts.FileProbe.Enabled {
		probe.WriteFile(opts.FileProbe.ReadinessPath)
	}
	log.Print("warm up finished")
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

	// Memory profile
	if opts.Profile.Memory != "" {

		log.Printf("memory profile will be written to %s file", opts.Profile.Memory)
		f, err := os.Create(opts.Profile.Memory)
		if err != nil {
			log.Fatalf("could not create memory profile: %v", err)
		}
		defer f.Close()

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
	}
}

func start(port int, livenessPath, readinessPath string, stop <-chan struct{}, done chan struct{}) *probe.Server {

	serverErr := make(chan struct{})
	probeServer := probe.NewServer(port, livenessPath, readinessPath)
	go func() {
		if err := probeServer.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Printf("probe server: %v", err)
				close(serverErr)
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-serverErr:
			log.Print("received probe server error")
			close(done)
		case <-stop:
			log.Print("received stop signal")
			probeServer.Shutdown()
			close(done)
		case sig := <-sigs:
			log.Printf("received %s signal", sig)
			probeServer.Shutdown()
			close(done)
		}
	}()
	return probeServer
}

// 'fan in' see: https://blog.golang.org/pipelines
func merge(cs ...<-chan warmup.Response) <-chan warmup.Response {

	var wg sync.WaitGroup
	out := make(chan warmup.Response)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan warmup.Response) {
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
