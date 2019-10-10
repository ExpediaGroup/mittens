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
	"mittens/pkg/warmup"
	"strings"
)

var allowedHttpMethods = map[string]interface{}{
	"GET":     nil,
	"HEAD":    nil,
	"POST":    nil,
	"PUT":     nil,
	"PATCH":   nil,
	"DELETE":  nil,
	"CONNECT": nil,
	"OPTIONS": nil,
	"TRACE":   nil,
}

type Http struct {
	Headers  []string
	Requests []string
}

func (h *Http) String() string {
	return fmt.Sprintf("%+v", *h)
}

func (h *Http) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(
		&h.Headers,
		"http-headers",
		nil,
		"Http headers to be sent with warm up requests.",
	)
	cmd.Flags().StringArrayVar(
		&h.Requests,
		"http-requests",
		nil,
		`Http request to be sent. Request is in '<http-method>:<path>[:body]' format. E.g. post:/ping:{"key": "value"}`,
	)
}

func (h *Http) GetWarmupHttpHeaders() map[string]string {
	return toHeaders(h.Headers)
}

func (h *Http) GetWarmupHttpRequests() ([]warmup.HttpRequest, error) {
	return toHttpRequests(h.Requests)
}

func toHttpRequests(requestsFlag []string) ([]warmup.HttpRequest, error) {

	var requests []warmup.HttpRequest
	for _, requestFlag := range requestsFlag {
		request, err := toHttpRequest(requestFlag)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func toHttpRequest(requestFlag string) (warmup.HttpRequest, error) {

	parts := strings.SplitN(requestFlag, ":", 3)
	if len(parts) < 2 {
		return warmup.HttpRequest{}, fmt.Errorf("invalid request flag: %s, expected format <http-method>:<path>[:body]", requestFlag)
	}

	method := strings.ToUpper(parts[0])
	_, ok := allowedHttpMethods[method]
	if !ok {
		return warmup.HttpRequest{}, fmt.Errorf("invalid request flag: %s, method %s is not supported", requestFlag, method)
	}

	// <method>:<path>
	if len(parts) == 2 {
		return warmup.HttpRequest{
			Method: method,
			Path:   parts[1],
			Body:   nil,
		}, nil
	}

	return warmup.HttpRequest{
		Method: method,
		Path:   parts[1],
		Body:   []byte(parts[2]),
	}, nil
}