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

// This is the main program which runs the warmup routine.
// Depending on the requests' types (HTTP or/and gRPC) it configures an HTTP client and/or grpcurl
// and sends requests for a fixed amount of time.

package routine

import (
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type StringArray []string

func (array *StringArray) String() string {
	var str strings.Builder

	for i := 0; i < len(*array); i++ {
		str.WriteString((*array)[i])
		str.WriteString("; ")
	}

	return str.String()
}

func (array *StringArray) Set(value string) error {
	*array = append(*array, value)
	return nil
}

type Options struct {
	ReadinessPath                   string
	TimeoutForReadinessProbeSeconds int
	TargetHost                      string
	TargetHttpPort                  int
	TargetGrpcPort                  int
	TargetProtocol                  string
	DurationSeconds                 int
	HttpTimeoutSeconds              int
	Concurrency                     int
	HttpHeaders                     StringArray
	GrpcHeaders                     StringArray
	WarmupRequests                  StringArray
	ExitAfterWarmup                 bool
	RequestDelayMilliseconds        int
}

func Warmup(options Options) {
	validateOptions(options)

	log.Printf("starting routine with arguments:\n")
	log.Printf("readinessPath: %s\n", options.ReadinessPath)
	log.Printf("timeoutForReadinessProbeSeconds: %v\n", options.TimeoutForReadinessProbeSeconds)
	log.Printf("targetHost: %v\n", options.TargetHost)
	log.Printf("targetProtocol: %v\n", options.TargetProtocol)
	log.Printf("targetHttpPort: %v\n", options.TargetHttpPort)
	log.Printf("targetGrpcPort: %v\n", options.TargetGrpcPort)
	log.Printf("durationSeconds: %v\n", options.DurationSeconds)
	log.Printf("httpTimeoutSeconds: %v\n", options.HttpTimeoutSeconds)
	log.Printf("concurrency: %v\n", options.Concurrency)
	log.Printf("exitAfterWarmup: %v\n", options.ExitAfterWarmup)
	log.Printf("httpHeaders: %v\n", options.HttpHeaders.String())
	log.Printf("grpcHeaders: %v\n", options.GrpcHeaders.String())
	log.Printf("warmupRequests: %v\n", options.WarmupRequests.String())
	log.Printf("requestDelay: %v\n", options.RequestDelayMilliseconds)

	// write a dummy file that can be used as a liveness check
	writeFile("alive")

	httpClient := createHttpClient(options.HttpTimeoutSeconds)

	// wait for readiness probe before running the warmup routine
	var timeElapsed time.Duration
	var err error

	if options.TimeoutForReadinessProbeSeconds == -1 {
		log.Printf("timeoutForReadinessProbeSeconds not set. will wait for readiness probe for a max of %d seconds\n", options.DurationSeconds)
		timeElapsed, err = waitForReadinessProbe(httpClient, options.TargetProtocol, options.TargetHost, options.TargetHttpPort, options.DurationSeconds, options.ReadinessPath)
	} else {
		timeElapsed, err = waitForReadinessProbe(httpClient, options.TargetProtocol, options.TargetHost, options.TargetHttpPort, options.TimeoutForReadinessProbeSeconds, options.ReadinessPath)
	}

	// run warmup only if container is ready. otherwise there's no point
	if err != nil {
		log.Printf("error: %s. will not proceed with warmup\n", err)
	} else {
		// if a timeout for readiness probe was not specified, deduct time already spent waiting
		var timeRemaining int

		if options.TimeoutForReadinessProbeSeconds == -1 {
			timeRemaining = options.DurationSeconds - int(timeElapsed/time.Second)
			log.Printf("timeoutForReadinessProbeSeconds not set. will run warmup for %d seconds\n", timeRemaining)
			if err = warmup(httpClient, options, timeRemaining); err != nil {
				log.Printf("warmup failed with error %v\n", err)
			}
		} else {
			log.Printf("timeoutForReadinessProbeSeconds set. will run warmup for full durationSeconds (%d seconds)\n", options.DurationSeconds)
			if err = warmup(httpClient, options, options.DurationSeconds); err != nil {
				log.Printf("warmup failed with error %v\n", err)
			}
		}
	}

	// we write to a file once routine finishes
	// this is a best effort routine so we always do this, regardless if routine succeeded or failed
	writeFile("ready")

	// we need to keep program alive or kubernetes will restart it
	if !options.ExitAfterWarmup {
		blockForever()
	}
}

func validateOptions(options Options) {
	var err error = nil

	if len(options.WarmupRequests) == 0 {
		err = fmt.Errorf("at least one warmup request must be supplied")
	}

	if options.TargetProtocol != "http" && options.TargetProtocol != "https" {
		err = fmt.Errorf("invalid protocol: %s", options.TargetProtocol)
	}

	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func waitForReadinessProbe(httpClient *http.Client, targetProtocol string, targetHost string, targetHttpPort int, durationSeconds int, readinessPath string) (time.Duration, error) {
	log.Printf("Waiting for a maximum of %v seconds until %s is ready\n", durationSeconds, readinessPath)

	var err error
	elapsed := time.Duration(0)

	start := time.Now()
	for {
		now := time.Now()
		elapsed = now.Sub(start)

		// timeDuration uses nano seconds
		var duration = time.Duration(durationSeconds) * 1000 * 1000 * 1000

		var emptyMap map[string]string
		url := fmt.Sprintf("%s://%s:%v%s", targetProtocol, targetHost, targetHttpPort, readinessPath)
		statusCode, err := sendHttpRequest(httpClient, "GET", url, &emptyMap, "")

		if err == nil && statusCode == 200 {
			break
		}

		if elapsed > duration {
			return elapsed, fmt.Errorf("readiness probe failed to respond after %v seconds", durationSeconds)
		}

		time.Sleep(time.Second)
	}

	return elapsed, err
}

func warmup(httpClient *http.Client, options Options, durationSeconds int) error {
	log.Printf("Warming up for %v seconds with %v threads\n", durationSeconds, options.Concurrency)
	var requestCounter uint64

	var channel = make(chan string, options.Concurrency) // channel will be used as a queue. we want it to have the same size as # of threads
	var waitGroup sync.WaitGroup

	// if any grpc requests, create grpc client
	var clientConnection *grpc.ClientConn
	var descriptorSource grpcurl.DescriptorSource
	var err error
	for _, request := range options.WarmupRequests {
		if strings.HasPrefix(request, "grpc") {
			clientConnection, descriptorSource, err = createGrpcClient(options.GrpcHeaders, options.TargetHost, options.TargetGrpcPort)
			if err != nil {
				return err
			}
			break
		}
	}

	// collect http headers into a map
	headersMap := make(map[string]string)
	for _, httpHeader := range options.HttpHeaders {
		entry := strings.Split(httpHeader, "=")
		headersMap[entry[0]] = entry[1]
	}

	// This starts a number (=concurrency) of goroutines that wait for something to do
	waitGroup.Add(options.Concurrency)
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			for {
				warmupRequest, ok := <-channel
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					waitGroup.Done()
					return
				}
				requestParts := strings.SplitN(warmupRequest, ":", 2)
				if requestParts[0] == "http" {
					if httpRequest, err := createHttpRequest(requestParts[1]); err != nil {
						panic(err)
					} else {
						callHttp(options, httpRequest, &requestCounter, httpClient, headersMap)
					}
				} else if requestParts[0] == "grpc" {
					if grpcRequest, err := createGrpcRequest(requestParts[1]); err != nil {
						panic(err)
					} else {
						callGrpc(descriptorSource, clientConnection, options, grpcRequest, &requestCounter, options.GrpcHeaders)
					}
				} else {
					fmt.Printf("unknown architecture: %s. use one of 'grpc' or 'http'\n", requestParts[0])
				}
			}
		}()
	}

	// add jobs to channel/queue
	start := time.Now()
	for {
		index := rand.Intn(len(options.WarmupRequests))
		warmupRequest := options.WarmupRequests[index]
		channel <- warmupRequest // add url to the queue
		t := time.Now()
		elapsed := t.Sub(start)

		// timeDuration uses nano seconds
		var duration = time.Duration(durationSeconds) * time.Second

		if elapsed > duration {
			log.Printf("stop adding new jobs to the queue\n")
			break
		}

		// add artificial delay to throttle requests
		time.Sleep(time.Duration(options.RequestDelayMilliseconds) * time.Millisecond)
	}

	log.Printf("warmup finished\n")

	close(channel)
	waitGroup.Wait() // wait for threads to finish

	log.Printf("Ran for %v seconds and made %v request(s)\n", time.Now().Sub(start).Seconds(), atomic.LoadUint64(&requestCounter))
	return nil
}

func callHttp(options Options, warmupRequest httpRequest, requestCounter *uint64, httpClient *http.Client, headersMap map[string]string) {
	url := fmt.Sprintf("%s://%s:%v%s", options.TargetProtocol, options.TargetHost, options.TargetHttpPort, warmupRequest.url)
	atomic.AddUint64(requestCounter, 1)
	statusCode, err := sendHttpRequest(httpClient, warmupRequest.method, url, &headersMap, warmupRequest.body)
	// if there was an error then sleep for a bit as not to overload the server
	if err != nil {
		log.Printf("Http call failed with error %v\n", err)
		time.Sleep(time.Second)
	} else if statusCode != 200 {
		log.Printf("Http call failed with status %d\n", statusCode)
		time.Sleep(time.Second)
	}
}

func callGrpc(descSource grpcurl.DescriptorSource, clientConnection *grpc.ClientConn, options Options, grpcRequest grpcRequest, requestCounter *uint64, grpcHeaders []string) {
	atomic.AddUint64(requestCounter, 1)
	err := sendGrpcRequest(descSource, clientConnection, grpcHeaders, grpcRequest)
	// if there was an error then sleep for a bit as not to overload the server
	if err != nil {
		log.Printf("gRPC call failed with error %v\n", err)
		time.Sleep(time.Second)
	}
}

func blockForever() {
	select {} // select blocks until one of its cases to return. there are no cases so it'll block indefinitely
}
