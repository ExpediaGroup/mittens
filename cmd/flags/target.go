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
	"fmt"
	"github.com/spf13/cobra"
	"mittens/pkg/grpc"
	"mittens/pkg/http"
	"mittens/pkg/warmup"
)

type Target struct {
	HttpHost               string
	HttpPort               int
	GrpcHost               string
	GrpcPort               int
	ReadinessPath          string
	ReadinessTimoutSeconds int
	Insecure               bool
}

func (t *Target) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t *Target) InitFlags(cmd *cobra.Command) {

	cmd.Flags().StringVar(
		&t.HttpHost,
		"target-http-host",
		"http://localhost",
		"Http host to warm up",
	)
	cmd.Flags().IntVar(
		&t.HttpPort,
		"target-http-port",
		8080,
		"Http port for warm up requests",
	)
	cmd.Flags().StringVar(
		&t.GrpcHost,
		"target-grpc-host",
		"localhost",
		"Grpc host to warm up",
	)
	cmd.Flags().IntVar(
		&t.GrpcPort,
		"target-grpc-port",
		50051,
		"Grpc port for warm up requests",
	)
	cmd.Flags().StringVar(
		&t.ReadinessPath,
		"target-readiness-path",
		"/ready",
		"The path used for target readiness probe",
	)
	cmd.Flags().IntVar(
		&t.ReadinessTimoutSeconds,
		"target-readiness-timeout-seconds",
		-1,
		"Timeout for target readiness probe",
	)
	cmd.Flags().BoolVar(
		&t.Insecure,
		"target-insecure",
		false,
		"Whether to skip TLS validation",
	)
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
