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
	HttpHost               string `json:"target-http-host"`
	HttpPort               int    `json:"target-http-port"`
	GrpcHost               string `json:"target-grpc-host"`
	GrpcPort               int    `json:"target-grpc-port"`
	ReadinessPath          string `json:"target-readiness-path"`
	ReadinessPort          int    `json:"target-readiness-port"`
	ReadinessTimoutSeconds int    `json:"target-readiness-timeout-seconds"`
	Insecure               bool   `json:"target-insecure"`
}

func (t *Target) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t *Target) InitFlags() {
	flag.StringVar(&t.HttpHost, "target-http-host", "http://localhost", "Http host to warm up")
	flag.IntVar(&t.HttpPort, "target-http-port", 8080, "Http port for warm up requests")
	flag.StringVar(&t.GrpcHost, "target-grpc-host", "localhost", "Grpc host to warm up")
	flag.IntVar(&t.GrpcPort, "target-grpc-port", 50051, "Grpc port for warm up requests")
	flag.StringVar(&t.ReadinessPath, "target-readiness-path", "/ready", "The path used for target readiness probe")
	flag.IntVar(&t.ReadinessPath, "target-readiness-port", "8080", "The port used for target readiness probe")
	flag.IntVar(&t.ReadinessTimoutSeconds, "target-readiness-timeout-seconds", -1, "Timeout for target readiness probe")
	flag.BoolVar(&t.Insecure, "target-insecure", false, "Whether to skip TLS validation")
}

func (t *Target) GetWarmupTargetOptions() warmup.TargetOptions {

	return warmup.TargetOptions{
		ReadinessPath:             t.ReadinessPath,
		ReadinessPort:             t.ReadinessPort,
		ReadinessTimeoutInSeconds: t.ReadinessTimoutSeconds,
	}
}

func (t *Target) GetReadinessHttpClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HttpHost, t.ReadinessPort), t.Insecure)
}

func (t *Target) GetHttpClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HttpHost, t.HttpPort), t.Insecure)
}

func (t *Target) GetGrpcClient(timeoutSeconds int) grpc.Client {
	return grpc.NewClient(fmt.Sprintf("%s:%d", t.GrpcHost, t.GrpcPort), t.Insecure, timeoutSeconds)
}
