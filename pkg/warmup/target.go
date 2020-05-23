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
	whttp "mittens/pkg/http"
	"net/http"
	"net/url"
	"time"
)

type TargetOptions struct {
	URL                       string
	ReadinessProtocol         string
	ReadinessHttpPath         string
	ReadinessGrpcMethod       string
	ReadinessPort             int
	ReadinessTimeoutInSeconds int
}

type Target struct {
	readinessHttpClient whttp.Client
	readinessGrpcClient grpc.Client
	httpClient          whttp.Client
	grpcClient          grpc.Client
	options             TargetOptions
}

func NewTarget(readinessHttpClient whttp.Client, readinessGrpcClient grpc.Client, httpClient whttp.Client, grpcClient grpc.Client, options TargetOptions) (Target,
	error) {

	if _, err := url.Parse(options.URL); err != nil {
		return Target{}, fmt.Errorf("invalid target host: %v", err)
	}

	t := Target{
		readinessHttpClient: readinessHttpClient,
		readinessGrpcClient: readinessGrpcClient,
		httpClient:          httpClient,
		grpcClient:          grpcClient,
		options:             options,
	}

	err := t.waitForReadinessProbe()
	if err != nil {
		return Target{}, fmt.Errorf("Error waiting for target to be ready: %v", err)
	}
	return t, err
}

func (t Target) waitForReadinessProbe() error {
	log.Printf("Waiting for target to be ready for a max of %ds", t.options.ReadinessTimeoutInSeconds)

	timeout := time.After(time.Duration(t.options.ReadinessTimeoutInSeconds) * time.Second)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("Giving up! Target not ready after %d seconds 🙁", t.options.ReadinessTimeoutInSeconds)
		default:
			// Wait one second between attempts. This is not configurable
			time.Sleep(time.Second * 1)

			if t.options.ReadinessProtocol == "http" {
				if err := t.readinessHttpClient.Request(http.MethodGet, t.options.ReadinessHttpPath, nil, nil); err.Err != nil {
					log.Printf("Target not ready yet...")
					continue
				}
			} else {
				request, err := grpc.ToGrpcRequest(t.options.ReadinessGrpcMethod)
				if err == nil {
					err1 := t.readinessGrpcClient.Request(request.ServiceMethod, "", nil)
					if err1.Err != nil {
						log.Printf("Target not ready yet...")
						continue
					}
				}
			}
			return nil
		}
	}
}
