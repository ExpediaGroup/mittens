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
	HttpHost               string `json:"httpHost"`
	HttpPort               int `json:"httpPort"`
	GrpcHost               string `json:"grpcHost"`
	GrpcPort               int `json:"grpcPort"`
	ReadinessPath          string `json:"readinessPath"`
	ReadinessTimoutSeconds int `json:"readinessTimeout"`
	Insecure               bool `json:"insecure"`
}

func (t *Target) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t *Target) InitFlags() {
	flag.StringVar(&t.HttpHost,"targetHttpHost","http://localhost","Http host to warm up")
	flag.IntVar(&t.HttpPort,"targetHttpPort",8080,"Http port for warm up requests")
	flag.StringVar(&t.GrpcHost,"targetGrpcHost","localhost","Grpc host to warm up")
	flag.IntVar(&t.GrpcPort,"targetGrpcPort",50051,"Grpc port for warm up requests")
	flag.StringVar(&t.ReadinessPath,"targetReadinessPath","/ready","The path used for target readiness probe")
	flag.IntVar(&t.ReadinessTimoutSeconds,"targetReadinessTimeoutSeconds",-1,"Timeout for target readiness probe")
	flag.BoolVar(&t.Insecure,"targetInsecure",false,"Whether to skip TLS validation")
}

func (t *Target) GetWarmupTargetOptions() warmup.TargetOptions {

	return warmup.TargetOptions{
		ReadinessPath:             t.ReadinessPath,
		ReadinessTimeoutInSeconds: t.ReadinessTimoutSeconds,
	}
}

func (t *Target) GetHttpClient() http.Client {
	return http.NewClient(fmt.Sprintf("%s:%d", t.HttpHost, t.HttpPort), t.Insecure)
}

func (t *Target) GetGrpcClient(timeoutSeconds int) grpc.Client {
	return grpc.NewClient(fmt.Sprintf("%s:%d", t.GrpcHost, t.GrpcPort), t.Insecure, timeoutSeconds)
}
