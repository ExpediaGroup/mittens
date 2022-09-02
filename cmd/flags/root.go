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
	"mittens/internal/pkg/grpc"
	"mittens/internal/pkg/http"
	"mittens/internal/pkg/warmup"
)

// Root stores all the flags.
type Root struct {
	MaxDurationSeconds       int
	MaxReadinessWaitSeconds  int
	MaxWarmupDurationSeconds int
	Concurrency              int
	RequestDelayMilliseconds int
	ConcurrencyTargetSeconds int
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
	// TODO: rename this to `max-global-duration-seconds`
	flag.IntVar(&r.MaxDurationSeconds, "max-duration-seconds", 60, "Global maximum duration. This includes both the time spent warming up the target service and also the time waiting for the target to become ready")
	flag.IntVar(&r.MaxReadinessWaitSeconds, "max-readiness-wait-seconds", 30, "Maximum time to wait for the target to become ready")
	flag.IntVar(&r.MaxWarmupDurationSeconds, "max-warmup-seconds", 30, "Maximum time spent sending warmup requests to the target service. Please note that `max-duration-seconds` may cap this duration.")
	flag.IntVar(&r.Concurrency, "concurrency", 2, "Number of concurrent requests for warm up")
	flag.IntVar(&r.RequestDelayMilliseconds, "request-delay-milliseconds", 500, "Delay in milliseconds between requests")
	flag.IntVar(&r.ConcurrencyTargetSeconds, "concurrency-target-seconds", 0, "Time taken to reach expected concurrency. This is useful to ramp up traffic.")
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

// GetMaxReadinessWaitSeconds returns the value of the max-readiness-wait-seconds parameter.
func (r *Root) GetMaxReadinessWaitSeconds() int {
	return r.MaxReadinessWaitSeconds
}

// GetMaxWarmupDurationSeconds returns the value of the max-warmup-seconds parameter.
func (r *Root) GetMaxWarmupDurationSeconds() int {
	return r.MaxWarmupDurationSeconds
}

// GetConcurrencyTargetSeconds returns the value of the concurrency-target-seconds parameter.
func (r *Root) GetConcurrencyTargetSeconds() int {
	return r.ConcurrencyTargetSeconds
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
	return r.Target.getReadinessGrpcClient()
}

// GetHTTPClient creates the HTTP client to be used for the actual requests.
func (r *Root) GetHTTPClient() http.Client {
	return r.Target.getHTTPClient()
}

// GetGrpcClient creates the gRPC client to be used for the actual requests.
func (r *Root) GetGrpcClient() grpc.Client {
	return r.Target.getGrpcClient()
}

// GetWarmupTargetOptions validates and returns any options that apply to the target.
func (r *Root) GetWarmupTargetOptions() (warmup.TargetOptions, error) {
	options := r.Target.getWarmupTargetOptions()
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

// GetWarmupHTTPRequests HTTP requests.
func (r *Root) GetWarmupHTTPRequests() ([]http.Request, error) {
	requests, err := r.HTTP.getWarmupHTTPRequests()
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetWarmupGrpcRequests returns gRPC requests.
func (r *Root) GetWarmupGrpcRequests() ([]grpc.Request, error) {
	requests, err := r.Grpc.getWarmupGrpcRequests()
	if err != nil {
		return nil, err
	}
	return requests, nil
}
