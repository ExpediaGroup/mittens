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
	"math/rand"
	"mittens/internal/pkg/grpc"
	"mittens/internal/pkg/http"
	"mittens/internal/pkg/safe"
	"mittens/internal/pkg/warmup"
	"time"
)

// Root stores all the flags.
type Root struct {
	MaxDurationSeconds       int
	Concurrency              int
	RequestDelayMilliseconds int
	ExitAfterWarmup          bool
	FailReadiness            bool
	FileProbe
	Target
	HTTP
	HTTPHeaders
	Grpc
}

func (r *Root) String() string {
	return fmt.Sprintf("%+v", *r)
}

// InitFlags initialises all the flags.
func (r *Root) InitFlags() {
	flag.IntVar(&r.MaxDurationSeconds, "max-duration-seconds", 60, "Max duration in seconds after which warm up will stop making requests")
	flag.IntVar(&r.Concurrency, "concurrency", 2, "Number of concurrent requests for warm up")
	flag.IntVar(&r.RequestDelayMilliseconds, "request-delay-milliseconds", 500, "Delay in milliseconds between requests")
	flag.BoolVar(&r.ExitAfterWarmup, "exit-after-warmup", false, "If warm up process should finish after completion. This is useful to prevent container restarts.")
	flag.BoolVar(&r.FailReadiness, "fail-readiness", false, "If set to true readiness will fail if no requests were sent.")

	r.FileProbe.initFlags()
	r.Target.initFlags()
	r.HTTPHeaders.initFlags()
	r.HTTP.initFlags()
	r.Grpc.initFlags()
}

// GetMaxDurationSeconds returns the value of the max-duration-seconds parameter.
func (r *Root) GetMaxDurationSeconds() int {
	return r.MaxDurationSeconds
}

// GetConcurrency returns the value of the concurrency parameter.
func (r *Root) GetConcurrency() int {
	return r.Concurrency
}

// GetReadinessHTTPClient creates the HTTP client to be used for the readiness requests.
func (r *Root) GetReadinessHTTPClient() http.Client {
	return r.Target.getReadinessHTTPClient()
}

// GetReadinessGrpcClient creates the gRPC client to be used for the readiness requests.
func (r *Root) GetReadinessGrpcClient() grpc.Client {
	return r.Target.getReadinessGrpcClient(r.MaxDurationSeconds)
}

// GetHTTPClient creates the HTTP client to be used for the actual requests.
func (r *Root) GetHTTPClient() http.Client {
	return r.Target.getHTTPClient()
}

// GetGrpcClient creates the gRPC client to be used for the actual requests.
func (r *Root) GetGrpcClient() grpc.Client {
	return r.Target.getGrpcClient(r.MaxDurationSeconds)
}

// GetWarmupTargetOptions validates and returns any options that apply to the target.
func (r *Root) GetWarmupTargetOptions() (warmup.TargetOptions, error) {
	options := r.Target.getWarmupTargetOptions()
	options.ReadinessTimeoutInSeconds = r.MaxDurationSeconds
	if options.ReadinessProtocol != "http" && options.ReadinessProtocol != "grpc" {
		err := fmt.Errorf("readiness protocol %s not supported, please use http or grpc", r.ReadinessProtocol)
		return options, err
	}
	return options, nil
}

// GetWarmupHTTPHeaders returns the HTTP headers.
func (r *Root) GetWarmupHTTPHeaders() []string {
	return r.HTTPHeaders.getWarmupHTTPHeaders()
}

// GetWarmupHTTPRequests returns a channel with HTTP requests.
func (r *Root) GetWarmupHTTPRequests() (chan http.Request, error) {
	requests, err := r.HTTP.getWarmupHTTPRequests()
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan http.Request)

	// create a goroutine that continuously adds requests to a channel for a maximum of MaxDurationSeconds
	go safe.Do(func() {
		if len(requests) == 0 {
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
	})
	return requestsChan, nil
}

// GetWarmupGrpcRequests returns a channel with gRPC requests.
func (r *Root) GetWarmupGrpcRequests() (chan grpc.Request, error) {
	requests, err := r.Grpc.getWarmupGrpcRequests()
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan grpc.Request)

	// create a goroutine that continuously adds requests to a channel for a maximum of MaxDurationSeconds
	go safe.Do(func() {
		if len(requests) == 0 {
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
	})
	return requestsChan, nil
}
