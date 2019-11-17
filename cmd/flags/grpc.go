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
	"log"
	"mittens/pkg/warmup"
	"strings"
)

type Grpc struct {
	Headers  stringArray `json:"grpc-headers"`
	Requests stringArray `json:"grpc-requests"`
}

func (g *Grpc) String() string {
	return fmt.Sprintf("%+v", *g)
}

func (g *Grpc) InitFlags() {
	flag.Var(&g.Headers, "grpc-headers", "gRPC header to be sent with warm up requests.")
	flag.Var(&g.Requests, "grpc-requests", `gRPC request to be sent. Request is in '<service>/<method>[:message]' format. E.g. health/ping:{"key": "value"}`)
}

func (g *Grpc) GetWarmupGrpcHeaders() []string {
	return g.Headers
}

func (g *Grpc) GetWarmupGrpcRequests() ([]warmup.GrpcRequest, error) {
	log.Print(g.Requests)
	return toGrpcRequests(g.Requests)
}

func toGrpcRequests(requestsFlag []string) ([]warmup.GrpcRequest, error) {

	var requests []warmup.GrpcRequest
	for _, requestFlag := range requestsFlag {
		log.Print(requestsFlag)
		request, err := toGrpcRequest(requestFlag)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func toGrpcRequest(requestFlag string) (warmup.GrpcRequest, error) {

	// service/method[:message]
	parts := strings.SplitN(requestFlag, ":", 2)
	log.Print(parts)
	if len(strings.Split(parts[0], "/")) != 2 {
		return warmup.GrpcRequest{}, fmt.Errorf("invalid request flag: %s, expected format <service>/<method>[:body]", requestFlag)
	}

	request := warmup.GrpcRequest{ServiceMethod: parts[0]}
	if len(parts) == 2 {
		request.Message = []byte(parts[1])
	}
	return request, nil
}
