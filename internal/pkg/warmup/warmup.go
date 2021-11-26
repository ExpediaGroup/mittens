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

package warmup

import (
	"log"
	"math/rand"
	"mittens/internal/pkg/grpc"
	"mittens/internal/pkg/http"
	"mittens/internal/pkg/safe"
	"mittens/internal/pkg/util"
	"sync"
	"time"
)

// Warmup holds any information needed for the workers to send requests.
type Warmup struct {
	Target                   Target
	MaxDurationSeconds       int
	Concurrency              int
	HttpRequests             chan http.Request
	HttpHeaders              []string
	GrpcRequests             chan grpc.Request
	RequestDelayMilliseconds int
	RampUpIntervalSeconds    int
}

// Run sends requests to the target using goroutines.
func (w Warmup) Run(hasHttpRequests bool, hasGrpcRequests bool, requestsSentCounter *int) {
	rand.Seed(time.Now().UnixNano()) // initialize seed only once to prevent deterministic/repeated calls every time we run

	var wg sync.WaitGroup

	if hasGrpcRequests {
		// connect to gRPC server once and only if there are gRPC requests
		log.Print("gRPC client connecting...")
		connErr := w.Target.grpcClient.Connect(w.HttpHeaders)

		if connErr != nil {
			log.Printf("gRPC client connect error: %v", connErr)
		} else {
			for i := 1; i <= w.Concurrency; i++ {
				delayIfNeeded(w.RampUpIntervalSeconds, i)
				log.Printf("Spawning new go routine for gRPC requests")
				wg.Add(1)
				go safe.Do(func() {
					w.GrpcWarmupWorker(&wg, w.GrpcRequests, w.HttpHeaders, w.RequestDelayMilliseconds, requestsSentCounter)
				})
			}
		}
	}

	if hasHttpRequests {
		for i := 1; i <= w.Concurrency; i++ {
			delayIfNeeded(w.RampUpIntervalSeconds, i)
			log.Printf("Spawning new go routine for HTTP requests")
			wg.Add(1)
			go safe.Do(func() {
				w.HTTPWarmupWorker(&wg, w.HttpRequests, util.ToHeaders(w.HttpHeaders), w.RequestDelayMilliseconds, requestsSentCounter)
			})
		}
	}

	wg.Wait()
}

// HTTPWarmupWorker sends HTTP requests to the target using goroutines.
func (w Warmup) HTTPWarmupWorker(wg *sync.WaitGroup, requests <-chan http.Request, headers map[string]string, requestDelayMilliseconds int, requestsSentCounter *int) {
	for request := range requests {
		time.Sleep(time.Duration(requestDelayMilliseconds) * time.Millisecond)

		resp := w.Target.httpClient.SendRequest(request.Method, request.Path, headers, request.Body)

		if resp.Err != nil {
			log.Printf("🔴 Error in request for %s: %v", request.Path, resp.Err)
		} else {
			*requestsSentCounter++

			if resp.StatusCode/100 == 2 {
				log.Printf("🟢 %s response\t%d ms\t%v\t%s\t%s", resp.Type, resp.Duration/time.Millisecond, resp.StatusCode, request.Method, request.Path)
			} else {
				log.Printf("🔴 %s response\t%d ms\t%v\t%s\t%s", resp.Type, resp.Duration/time.Millisecond, resp.StatusCode, request.Method, request.Path)
			}
		}
	}
	wg.Done()
}

// GrpcWarmupWorker sends gRPC requests to the target using goroutines.
func (w Warmup) GrpcWarmupWorker(wg *sync.WaitGroup, requests <-chan grpc.Request, headers []string, requestDelayMilliseconds int, requestsSentCounter *int) {
	for request := range requests {
		time.Sleep(time.Duration(requestDelayMilliseconds) * time.Millisecond)

		resp := w.Target.grpcClient.SendRequest(request.ServiceMethod, request.Message, headers, false)

		if resp.Err != nil {
			log.Printf("🔴 Error in request for %s: %v", request.ServiceMethod, resp.Err)
		} else {
			*requestsSentCounter++
			log.Printf("🟢 %s response\t%d ms %s", resp.Type, resp.Duration/time.Millisecond, request.ServiceMethod)
		}

	}
	wg.Done()
}

func delayIfNeeded(rampUpIntervalSeconds int, currentConcurrency int) {
	if rampUpIntervalSeconds > 0 && currentConcurrency > 1 {
		time.Sleep(time.Duration(rampUpIntervalSeconds) * time.Second)
	}
}
