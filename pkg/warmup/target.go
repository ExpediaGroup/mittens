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
	ReadinessPath             string
	ReadinessTimeoutInSeconds int
}

type Target struct {
	httpClient whttp.Client
	grpcClient grpc.Client
	options    TargetOptions
}

func NewTarget(httpClient whttp.Client, grpcClient grpc.Client, options TargetOptions, done <-chan struct{}) (Target, error) {

	if _, err := url.Parse(options.URL); err != nil {
		return Target{}, fmt.Errorf("invalid target host: %v", err)
	}

	t := Target{
		httpClient: httpClient,
		grpcClient: grpcClient,
		options:    options,
	}

	if err := t.waitForReadinessProbe(done); err != nil {
		return Target{}, fmt.Errorf("wait for rediness probe: %v", err)
	}
	return t, nil
}

func (t Target) waitForReadinessProbe(done <-chan struct{}) error {

	log.Printf("waiting for %s target readiness probe", t.options.ReadinessPath)
	for {
		select {
		case <-time.After(time.Duration(t.options.ReadinessTimeoutInSeconds) * time.Second):
			return fmt.Errorf("target readiness proble: timeout %d seconds exceeded", t.options.ReadinessTimeoutInSeconds)
		case <-done:
			return errors.New("target readiness probe: received done signal")
		default:
			if err := t.httpClient.Request(http.MethodGet, t.options.ReadinessPath, nil, nil); err != nil {
				continue
				log.Printf("target readiness probe: %v", err)
			}
			log.Print("target is ready")
			return nil
		}
	}
	return nil
}