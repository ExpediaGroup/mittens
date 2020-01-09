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
	"errors"
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
	readinessGrpcClient	grpc.Client
	httpClient          whttp.Client
	grpcClient          grpc.Client
	options             TargetOptions
}

func NewTarget(readinessHttpClient whttp.Client, readinessGrpcClient grpc.Client, httpClient whttp.Client, grpcClient grpc.Client, options TargetOptions,
	done <-chan struct{}) (Target,
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

	err := t.waitForReadinessProbe(done)
	if err != nil {
		return Target{}, fmt.Errorf("wait for readiness probe: %v", err)
	}
	return t, err
}

func (t Target) waitForReadinessProbe(done <-chan struct{}) error {

	log.Printf("waiting for target readiness probe for %ds", t.options.ReadinessTimeoutInSeconds)

	timeout := time.After(time.Duration(t.options.ReadinessTimeoutInSeconds) * time.Second)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("target readiness probe: timeout %d seconds exceeded", t.options.ReadinessTimeoutInSeconds)
		case <-done:
			return errors.New("target readiness probe: received done signal")
		default:
			if t.options.ReadinessProtocol == "http" {
				if err := t.readinessHttpClient.Request(http.MethodGet, t.options.ReadinessHttpPath, nil, nil); err != nil {
					log.Printf("target readiness probe: %v", err)
					continue
				}
			} else {
				request, err := grpc.ToGrpcRequest(t.options.ReadinessGrpcMethod)
				if err == nil {
					err1 := t.readinessGrpcClient.Request(request.ServiceMethod, "", nil)
					if err1 != nil {
						log.Printf("target readiness probe: %v", err1)
						continue
					}
				}
			}
			log.Print("target is ready")
			return nil
		}
	}
}
