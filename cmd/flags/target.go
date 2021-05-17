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

// Target stores flags related to the target.
type Target struct {
	HTTPHost            string
	HTTPPort            int
	GrpcHost            string
	GrpcPort            int
	ReadinessProtocol   string
	ReadinessHTTPPath   string
	ReadinessGrpcMethod string
	ReadinessPort       int
	Insecure            bool
}

func (t *Target) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t *Target) initFlags() {
	flag.StringVar(&t.HTTPHost, "target-http-host", "http://localhost", "HTTP host to warm up")
	flag.IntVar(&t.HTTPPort, "target-http-port", 8080, "HTTP port for warm up requests")
	flag.StringVar(&t.GrpcHost, "target-grpc-host", "localhost", "Grpc host to warm up")
	flag.IntVar(&t.GrpcPort, "target-grpc-port", 50051, "Grpc port for warm up requests")
	flag.StringVar(&t.ReadinessProtocol, "target-readiness-protocol", "http", "Protocol to be used for readiness check. One of [http, grpc]")
	flag.StringVar(&t.ReadinessHTTPPath, "target-readiness-http-path", "/ready", "The path used for HTTP target readiness probe")
	flag.StringVar(&t.ReadinessGrpcMethod, "target-readiness-grpc-method", "grpc.health.v1.Health/Check", "The service method used for gRPC target readiness probe")
	flag.IntVar(&t.ReadinessPort, "target-readiness-port", toIntOrDefaultIfNull(&t.HTTPPort, 8080), "The port used for target readiness probe")
	flag.BoolVar(&t.Insecure, "target-insecure", false, "Whether to skip TLS validation")
}

func toIntOrDefaultIfNull(value *int, defaultValue int) int {
	if value == nil {
		return defaultValue
	}
	return *value
}

func (t *Target) getWarmupTargetOptions() warmup.TargetOptions {

	return warmup.TargetOptions{
		ReadinessProtocol:   t.ReadinessProtocol,
		ReadinessHTTPPath:   t.ReadinessHTTPPath,
		ReadinessGrpcMethod: t.ReadinessGrpcMethod,
		ReadinessPort:       t.ReadinessPort,
	}
}

func (t *Target) getReadinessHTTPClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HTTPHost, t.ReadinessPort), t.Insecure)
}

func (t *Target) getReadinessGrpcClient(timeoutSeconds int) grpc.Client {
	return grpc.NewClient(fmt.Sprintf("%s:%d", t.GrpcHost, t.ReadinessPort), t.Insecure, timeoutSeconds)
}

func (t *Target) getHTTPClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HTTPHost, t.HTTPPort), t.Insecure)
}

func (t *Target) getGrpcClient(timeoutSeconds int) grpc.Client {
	return grpc.NewClient(fmt.Sprintf("%s:%d", t.GrpcHost, t.GrpcPort), t.Insecure, timeoutSeconds)
}
