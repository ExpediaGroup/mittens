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

package routine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func TestWarmup(t *testing.T) {
	var ready uint64 = 0
	var get uint64 = 0
	var post uint64 = 0
	var put uint64 = 0

	postPayload := "{\"foo\":\"bar\"}"
	putPayload := "{\"baz\":\"qux\"}"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string

		if r.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			body = string(bodyBytes)
		}

		if r.Method == "GET" && r.URL.Path == "/ready" {
			atomic.AddUint64(&ready, 1)
		} else if r.Method == "GET" && r.URL.String() == "/getEndpoint" {
			atomic.AddUint64(&get, 1)
		} else if r.Method == "POST" && r.URL.String() == "/postEndpoint" && body == postPayload {
			atomic.AddUint64(&post, 1)
		} else if r.Method == "PUT" && r.URL.String() == "/putEndpoint" && body == putPayload {
			atomic.AddUint64(&put, 1)
		}
	}))
	defer server.Close()

	serverPort, _ := strconv.Atoi(strings.Split(strings.TrimPrefix(server.URL, "http://"), ":")[1])

	options := Options{
		ReadinessPath:                   "/ready",
		TimeoutForReadinessProbeSeconds: 3,
		TargetHttpPort:                  serverPort,
		TargetHost:                      "localhost",
		TargetProtocol:                  "http",
		DurationSeconds:                 1,
		HttpTimeoutSeconds:              3,
		Concurrency:                     1,
		HttpHeaders:                     []string{},
		ExitAfterWarmup:                 true,
		WarmupRequests: []string{
			"http:get:/getEndpoint",
			"http:get:/getEndpoint:",
			fmt.Sprintf("http:post:/postEndpoint:%s", postPayload),
			fmt.Sprintf("http:put:/putEndpoint:%s", putPayload),
		},
	}

	Warmup(options)

	if ready == 0 {
		t.Errorf("Ready endpoint not hit")
	}
	if get == 0 {
		t.Errorf("Get endpoint not hit")
	}
	if post == 0 {
		t.Errorf("Post endpoint not hit")
	}
	if put == 0 {
		t.Errorf("Put endpoint not hit")
	}
}
