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

package flags

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"mittens/pkg/grpc"
	"mittens/pkg/http"
	"mittens/pkg/warmup"
	"time"
)

type Root struct {
	MaxDurationSeconds       int  `json:"max-duration-seconds"`
	Concurrency              int  `json:"concurrency"`
	RequestDelayMilliseconds int  `json:"request-delay-milliseconds"`
	ExitAfterWarmup          bool `json:"exit-after-warmup"`
	FailReadiness            bool `json:"fail-readiness"`
	FileProbe
	ServerProbe
	Target
	Http
	Grpc
}

func (r *Root) String() string {
	return fmt.Sprintf("%+v", *r)
}

func (r *Root) InitFlags() {
	flag.IntVar(&r.MaxDurationSeconds, "max-duration-seconds", 60, "Max duration in seconds after which warm up will stop making requests")
	flag.IntVar(&r.Concurrency, "concurrency", 2, "Number of concurrent requests for warm up")
	flag.IntVar(&r.RequestDelayMilliseconds, "request-delay-milliseconds", 250, "Delay in milliseconds between requests")
	flag.BoolVar(&r.ExitAfterWarmup, "exit-after-warmup", false, "If warm up process should finish after completion. This is useful to prevent container restarts.")
	flag.BoolVar(&r.FailReadiness, "fail-readiness", false, "If set to true readiness will fail if no requests were sent.")

	r.FileProbe.InitFlags()
	r.ServerProbe.InitFlags()
	r.Target.InitFlags()
	r.Http.InitFlags()
	r.Grpc.InitFlags()
}

func (r *Root) GetWarmupOptions() warmup.Options {

	return warmup.Options{
		MaxDurationSeconds: r.MaxDurationSeconds,
		Concurrency:        r.Concurrency,
	}
}

func (r *Root) GetReadinessHttpClient() http.Client {
	return r.Target.GetReadinessHttpClient()
}

func (r *Root) GetReadinessGrpcClient() grpc.Client {
	return r.Target.GetReadinessGrpcClient()
}

func (r *Root) GetHttpClient() http.Client {
	return r.Target.GetHttpClient()
}

func (r *Root) GetGrpcClient() grpc.Client {
	return r.Target.GetGrpcClient(r.MaxDurationSeconds)
}

func (r *Root) GetWarmupTargetOptions() (warmup.TargetOptions, error) {

	options := r.Target.GetWarmupTargetOptions()
	if options.ReadinessTimeoutInSeconds <= 0 {
		log.Printf("Readiness timeout in seconds not set, defaulting to max duration in seconds: %ds", r.MaxDurationSeconds)
		options.ReadinessTimeoutInSeconds = r.MaxDurationSeconds
	}
	if options.ReadinessProtocol != "http" && options.ReadinessProtocol != "grpc" {
		err := fmt.Errorf("Readiness protocol %s not supported, please use http or grpc", r.ReadinessProtocol)
		return options, err
	}
	return options, nil
}

func (r *Root) GetWarmupHttpHeaders() map[string]string {
	return r.Http.GetWarmupHttpHeaders()
}

func (r *Root) GetWarmupHttpRequests() (chan http.Request, error) {
	requests, err := r.Http.GetWarmupHttpRequests()
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan http.Request)

	// create a goroutine that continuously adds requests to a channel for a maximum of MaxDurationSeconds
	go func() {
		if len(requests) == 0 {
			log.Print("No http warm up requests specified")
			close(requestsChan)
			return
		}
		timeout := time.After(time.Duration(r.MaxDurationSeconds) * time.Second)

		for {
			select {
			case <-timeout:
				close(requestsChan)
				return
			default:
				number := rand.Intn(len(requests))
				requestsChan <- requests[number]
			}
		}
	}()
	return requestsChan, nil
}

func (r *Root) GetWarmupGrpcRequests() (chan grpc.Request, error) {
	requests, err := r.Grpc.GetWarmupGrpcRequests()
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan grpc.Request)

	// create a goroutine that continuously adds requests to a channel for a maximum of MaxDurationSeconds
	go func() {
		if len(requests) == 0 {
			log.Print("No grpc warm up requests specified")
			close(requestsChan)
			return
		}
		timeout := time.After(time.Duration(r.MaxDurationSeconds) * time.Second)

		for {
			select {
			case <-timeout:
				close(requestsChan)
				return
			default:
				number := rand.Intn(len(requests))
				requestsChan <- requests[number]
			}
		}
	}()
	return requestsChan, nil
}

func (r *Root) GetWarmupGrpcHeaders() []string {
	return r.Grpc.GetWarmupGrpcHeaders()
}
