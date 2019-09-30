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

package main

import (
	"flag"
	"mittens/routine"
)

func main() {
	var options routine.Options

	flag.StringVar(&options.ReadinessPath, "readinessPath", "/ready", "the path used for http readiness probe")
	flag.IntVar(&options.TimeoutForReadinessProbeSeconds, "timeoutForReadinessProbeSeconds", -1, "how long to wait for readiness probe of target server")
	flag.StringVar(&options.TargetHost, "targetHost", "localhost", "the target server")
	flag.IntVar(&options.TargetHttpPort, "targetHttpPort", 8080, "the http port for the target server")
	flag.IntVar(&options.TargetGrpcPort, "targetGrpcPort", 50051, "the grpc port for the target server")
	flag.StringVar(&options.TargetProtocol, "targetProtocol", "http", "the protocol for the target server")
	flag.IntVar(&options.DurationSeconds, "durationSeconds", 60, "duration of the routine")
	flag.IntVar(&options.HttpTimeoutSeconds, "httpTimeoutSeconds", 10, "http timeout in seconds")
	flag.IntVar(&options.Concurrency, "concurrency", 2, "concurrent requests for routine")
	flag.BoolVar(&options.ExitAfterWarmup, "exitAfterWarmup", false, "exit after routine has run")
	flag.Var(&options.HttpHeaders, "httpHeader", "collection of http headers to be sent with warmup requests.")
	flag.Var(&options.GrpcHeaders, "grpcHeader", "collection of grpc headers to be sent with warmup requests.")
	flag.Var(&options.WarmupRequests, "warmupRequest", "collection of (relative) urls to be used during the warmup routine. {today} and {tomorrow} can be used to get today's and tomorrow's date in YYYY-MM-DD format")
	flag.IntVar(&options.RequestDelayMilliseconds, "requestDelayMilliseconds", 0, "adds a delay between requests")

	flag.Parse()

	routine.Warmup(options)
}
