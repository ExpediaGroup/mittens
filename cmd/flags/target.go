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
	"mittens/pkg/grpc"
	"mittens/pkg/http"
	"mittens/pkg/warmup"
)

type Target struct {
	HttpHost                string
	HttpPort                int
	GrpcHost                string
	GrpcPort                int
	ReadinessProtocol       string
	ReadinessHttpPath       string
	ReadinessGrpcMethod     string
	ReadinessPort           int
	ReadinessTimeoutSeconds int
	Insecure                bool
}

func (t *Target) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t *Target) InitFlags() {
	flag.StringVar(&t.HttpHost, "target-http-host", "http://localhost", "Http host to warm up")
	flag.IntVar(&t.HttpPort, "target-http-port", 8080, "Http port for warm up requests")
	flag.StringVar(&t.GrpcHost, "target-grpc-host", "localhost", "Grpc host to warm up")
	flag.IntVar(&t.GrpcPort, "target-grpc-port", 50051, "Grpc port for warm up requests")
	flag.StringVar(&t.ReadinessProtocol, "target-readiness-protocol", "http", "Protocol to be used for readiness check. One of [http, grpc]")
	flag.StringVar(&t.ReadinessHttpPath, "target-readiness-http-path", "/ready", "The path used for HTTP target readiness probe")
	flag.StringVar(&t.ReadinessGrpcMethod, "target-readiness-grpc-method", "grpc.health.v1.Health/Check", "The service method used for gRPC target readiness probe")
	flag.IntVar(&t.ReadinessPort, "target-readiness-port", toIntOrDefaultIfNull(&t.HttpPort, 8080), "The port used for target readiness probe")
	flag.BoolVar(&t.Insecure, "target-insecure", false, "Whether to skip TLS validation")
}

func toIntOrDefaultIfNull(value *int, defaultValue int) int {
	if value == nil {
		return defaultValue
	} else {
		return *value
	}
}

func (t *Target) GetWarmupTargetOptions() warmup.TargetOptions {

	return warmup.TargetOptions{
		ReadinessProtocol:         t.ReadinessProtocol,
		ReadinessHttpPath:         t.ReadinessHttpPath,
		ReadinessGrpcMethod:       t.ReadinessGrpcMethod,
		ReadinessPort:             t.ReadinessPort,
		ReadinessTimeoutInSeconds: t.ReadinessTimeoutSeconds,
	}
}

func (t *Target) GetReadinessHttpClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HttpHost, t.ReadinessPort), t.Insecure)
}

func (t *Target) GetReadinessGrpcClient() grpc.Client {
	return grpc.NewClient(fmt.Sprintf("%s:%d", t.GrpcHost, t.ReadinessPort), t.Insecure, t.ReadinessTimeoutSeconds)
}

func (t *Target) GetHttpClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HttpHost, t.HttpPort), t.Insecure)
}

func (t *Target) GetGrpcClient(timeoutSeconds int) grpc.Client {
	return grpc.NewClient(fmt.Sprintf("%s:%d", t.GrpcHost, t.GrpcPort), t.Insecure, timeoutSeconds)
}
