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
	"fmt"
	"log"
	"mittens/pkg/grpc"
	"mittens/pkg/http"
	"sync"
	"time"
)

type Options struct {
	MaxDurationSeconds int
	Concurrency        int
}

type Warmup struct {
	target  Target
	options Options
}

func NewWarmup(readinessHttpClient http.Client, readinessGrpcClient grpc.Client, httpClient http.Client, grpcClient grpc.Client, options Options, targetOptions TargetOptions) (Warmup, error) {
	target, err := NewTarget(readinessHttpClient, readinessGrpcClient, httpClient, grpcClient, targetOptions)
	if err != nil {
		return Warmup{}, fmt.Errorf("new target: %v", err)
	}
	return Warmup{target: target, options: options}, err
}

func (w Warmup) HttpWarmupWorker(wg *sync.WaitGroup, requests <-chan http.Request, headers map[string]string, requestDelayMilliseconds int, requestsSentCounter *int) {
	for request := range requests {
		time.Sleep(time.Duration(requestDelayMilliseconds) * time.Millisecond)

		resp := w.target.httpClient.Request(request.Method, request.Path, headers, request.Body)

		if resp.Err != nil {
			log.Printf("🔴 Error in request for %s: %v", request.Path, resp.Err)
		} else {
			*requestsSentCounter++

			if resp.StatusCode/100 == 2 {
				log.Printf("🟢 %s response for %s %d ms: %v", resp.Type, request.Path, resp.Duration/time.Millisecond, resp.StatusCode)
			} else {
				log.Printf("🔴 %s response for %s %d ms: %v", resp.Type, request.Path, resp.Duration/time.Millisecond, resp.StatusCode)
			}
		}
	}
	wg.Done()
}

func (w Warmup) GrpcWarmupWorker(wg *sync.WaitGroup, headers []string, requests <-chan grpc.Request, requestDelayMilliseconds int, requestsSentCounter *int) {
	for request := range requests {
		time.Sleep(time.Duration(requestDelayMilliseconds) * time.Millisecond)

		resp := w.target.grpcClient.Request(request.ServiceMethod, request.Message, headers)

		if resp.Err != nil {
			log.Printf("🔴 Error in request for %s: %v", request.ServiceMethod, resp.Err)
		} else {
			*requestsSentCounter++

			log.Printf("%s response for %s %d ms", resp.Type, request.ServiceMethod, resp.Duration/time.Millisecond)
		}

	}
	wg.Done()
}
