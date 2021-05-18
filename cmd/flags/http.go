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
	"mittens/internal/pkg/http"
)

var allowedHTTPMethods = map[string]interface{}{
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

// HTTP stores flags related to HTTP requests.
type HTTP struct {
	Requests stringArray
}

func (h *HTTP) String() string {
	return fmt.Sprintf("%+v", *h)
}

func (h *HTTP) initFlags() {
	flag.Var(&h.Requests, "http-requests", `HTTP request to be sent. Request is in '<http-method>:<path>[:body]' format. E.g. post:/ping:{"key":"value"}`)
}

func (h *HTTP) getWarmupHTTPRequests() ([]http.Request, error) {
	return toHTTPRequests(h.Requests)
}

func toHTTPRequests(requestsFlag []string) ([]http.Request, error) {
	var requests []http.Request
	for _, requestFlag := range requestsFlag {
		request, err := http.ToHTTPRequest(requestFlag)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}
