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
	"mittens/pkg/grpc"
)

// Grpc stores flags related to gRPC requests.
type Grpc struct {
	Headers  stringArray
	Requests stringArray
}

func (g *Grpc) String() string {
	return fmt.Sprintf("%+v", *g)
}

func (g *Grpc) initFlags() {
	flag.Var(&g.Headers, "grpc-headers", "gRPC header to be sent with warm up requests.")
	flag.Var(&g.Requests, "grpc-requests", `gRPC request to be sent. Request is in '<service>/<method>[:message]' format. E.g. health/ping:{"key": "value"}`)
}

func (g *Grpc) getWarmupGrpcHeaders() []string {
	return g.Headers
}

func (g *Grpc) getWarmupGrpcRequests() ([]grpc.Request, error) {
	log.Print(g.Requests)
	return toGrpcRequests(g.Requests)
}

func toGrpcRequests(requestsFlag []string) ([]grpc.Request, error) {

	var requests []grpc.Request
	for _, requestFlag := range requestsFlag {
		request, err := grpc.ToGrpcRequest(requestFlag)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}
