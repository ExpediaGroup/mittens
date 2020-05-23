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

func (w Warmup) spawnHttpWarmupWorker(wg *sync.WaitGroup, requests <-chan http.Request, headers map[string]string, requestDelayMilliseconds int) {
	for request := range requests {
		time.Sleep(time.Duration(requestDelayMilliseconds) * time.Millisecond)

		go func(r http.Request) {
			wg.Add(1)
			defer wg.Done()
			resp := w.target.httpClient.Request(r.Method, r.Path, headers, r.Body)

			if resp.Err != nil {
				log.Printf("🔴 %s response %d milliseconds: error: %v", resp.Type, resp.Duration/time.Millisecond, resp.Err)
			} else {
				log.Printf("🟢 %s response %d milliseconds: OK", resp.Type, resp.Duration/time.Millisecond)
			}
		}(request)
	}
	wg.Done()
}

func (w Warmup) HttpWarmup(headers map[string]string, requests <-chan http.Request, requestDelayMilliseconds int, concurrency int) {
	var wg sync.WaitGroup

	for i := 1; i <= concurrency; i++ {
		log.Printf("Spawning new go routine")
		wg.Add(1)
		go w.spawnHttpWarmupWorker(&wg, requests, headers, requestDelayMilliseconds)
	}

	wg.Wait()
}

func (w Warmup) spawnGrpcWarmupWorker(wg *sync.WaitGroup, headers []string, requests <-chan grpc.Request, requestDelayMilliseconds int) {
	for request := range requests {
		time.Sleep(time.Duration(requestDelayMilliseconds) * time.Millisecond)

		go func(r grpc.Request) {
			wg.Add(1)
			defer wg.Done()
			resp := w.target.grpcClient.Request(r.ServiceMethod, r.Message, headers)

			if resp.Err != nil {
				log.Printf("🔴 %s response %d milliseconds: error: %v", resp.Type, resp.Duration/time.Millisecond, resp.Err)
			} else {
				log.Printf("🟢 %s response %d milliseconds: OK", resp.Type, resp.Duration/time.Millisecond)
			}
		}(request)
	}
	wg.Done()
}

func (w Warmup) GrpcWarmup(headers []string, requests <-chan grpc.Request, requestDelayMilliseconds int, concurrency int) {
	var wg sync.WaitGroup

	for i := 1; i <= concurrency; i++ {
		log.Printf("Spawning new go routine")
		wg.Add(1)
		go w.spawnGrpcWarmupWorker(&wg, headers, requests, requestDelayMilliseconds)
	}

	wg.Wait()
}
