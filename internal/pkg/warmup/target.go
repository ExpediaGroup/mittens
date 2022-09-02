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
	"mittens/internal/pkg/grpc"
	whttp "mittens/internal/pkg/http"
	"net/http"
	"time"
)

// TargetOptions represents target configurations set by the user.
type TargetOptions struct {
	ReadinessProtocol   string
	ReadinessHTTPPath   string
	ReadinessGrpcMethod string
	ReadinessPort       int
}

// Target includes information needed to send requests to the target. It includes configured http and gRPC clients and options set by the user.
type Target struct {
	readinessHTTPClient whttp.Client
	readinessGrpcClient grpc.Client
	httpClient          whttp.Client
	grpcClient          grpc.Client
	options             TargetOptions
}

// NewTarget returns an instance of the target versus which mittens will run.
func NewTarget(readinessHTTPClient whttp.Client, readinessGrpcClient grpc.Client, httpClient whttp.Client, grpcClient grpc.Client, options TargetOptions) Target {
	t := Target{
		readinessHTTPClient: readinessHTTPClient,
		readinessGrpcClient: readinessGrpcClient,
		httpClient:          httpClient,
		grpcClient:          grpcClient,
		options:             options,
	}
	return t
}

// WaitForReadinessProbe sends health-check requests to the target and waits until it becomes ready.
// It returns an error if the timeout is exceeded.
// It supports both HTTP and gRPC health-checks.
func (t Target) WaitForReadinessProbe(maxReadinessWaitDurationInSeconds int, headers []string) error {
	log.Printf("Waiting for %s target to be ready for a max of %ds", t.options.ReadinessProtocol, maxReadinessWaitDurationInSeconds)

	timeout := time.After(time.Duration(maxReadinessWaitDurationInSeconds) * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("giving up; target not ready after %d seconds 🙁", maxReadinessWaitDurationInSeconds)
		default:
			// Wait one second between attempts. This is not configurable
			time.Sleep(time.Second * 1)

			if t.options.ReadinessProtocol == "http" {
				// error if error in the response or status code not in the 200 range
				if resp := t.readinessHTTPClient.SendRequest(http.MethodGet, t.options.ReadinessHTTPPath, headers, nil); resp.Err != nil || resp.StatusCode/100 != 2 {
					log.Printf("HTTP target not ready yet...")
					continue
				}
			} else {
				request, err := grpc.ToGrpcRequest(t.options.ReadinessGrpcMethod)
				if err == nil {
					log.Print("gRPC readiness client connecting...")
					connErr := t.readinessGrpcClient.Connect(nil)
					if connErr != nil {
						log.Printf("gRPC readiness client connect error: %v", connErr)
					}
					err1 := t.readinessGrpcClient.SendRequest(request.ServiceMethod, "", headers, false)
					if err1.Err != nil {
						log.Printf("gRPC target not ready yet...")
						continue
					}
				}
			}
			return nil
		}
	}
}
