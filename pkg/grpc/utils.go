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

package grpc

import (
	"fmt"
	"mittens/pkg/placeholders"
	"strings"
)

// Request represents a gRPC request.
type Request struct {
	ServiceMethod string
	Message       string
}

// ToGrpcRequest parses a gRPC request which is in a string format and stores it in a struct.
func ToGrpcRequest(requestFlag string) (Request, error) {

	// service/method[:message]
	parts := strings.SplitN(requestFlag, ":", 2)
	if len(strings.Split(parts[0], "/")) != 2 {
		return Request{}, fmt.Errorf("invalid request flag: %s, expected format <service>/<method>[:body]", requestFlag)
	}

	request := Request{ServiceMethod: parts[0]}
	if len(parts) == 2 {
		request.Message = placeholders.InterpolatePlaceholders(parts[1])
	} else {
		request.Message = ""
	}
	return request, nil
}
