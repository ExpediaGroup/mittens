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
	"mittens/pkg/warmup"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var todayTemplateRegex = regexp.MustCompile("{(today([+-]\\d+)?|tomorrow)}")
var todayTemplatePlusMinusRegex = regexp.MustCompile("[+-]\\d+")

type Http struct {
	Headers  stringArray `json:"http-headers"`
	Requests stringArray `json:"http-requests"`
}

func (h *Http) String() string {
	return fmt.Sprintf("%+v", *h)
}

func (h *Http) InitFlags() {
	flag.Var(&h.Headers, "http-headers", "Http header to be sent with warm up requests.")
	flag.Var(&h.Requests, "http-requests", `Http request to be sent. Request is in '<http-method>:<path>[:body]' format. E.g. post:/ping:{"key":"value"}`)
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
		path := interpolateDates(parts[1])

		return warmup.HttpRequest{
			Method: method,
			Path:   path,
			Body:   nil,
		}, nil
	}

	path := interpolateDates(parts[1])
	var body = interpolateDates(parts[2])

	return warmup.HttpRequest{
		Method: method,
		Path:   path,
		Body:   &body,
	}, nil
}

func interpolateDates(source string) string {
	return todayTemplateRegex.ReplaceAllStringFunc(source, func(templateString string) string {
		offsetDays := 0

		if templateString == "{tomorrow}" {
			offsetDays = 1
		} else if extractedOffset := todayTemplatePlusMinusRegex.FindString(templateString); len(extractedOffset) > 0 {
			offsetDays, _ = strconv.Atoi(extractedOffset)
		}

		// the date below is how the golang date formatter works. it's used for the formatting. it's not what is actually going to be displayed
		return time.Now().Add(time.Duration(offsetDays) * 24 * time.Hour).Format("2006-01-02")
	})
}
