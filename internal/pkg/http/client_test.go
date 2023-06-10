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

package http

import (
	"context"
	"fmt"
	"mittens/fixture"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockServer *http.Server

const WorkingPath = "/path"

var serverUrl string

func TestMain(m *testing.M) {
	setup()

	m.Run()
	teardown()
}

func TestRequestSuccess(t *testing.T) {
	c := NewClient(serverUrl, false, 10000)
	reqBody := ""
	resp := c.SendRequest("GET", WorkingPath, []string{}, &reqBody)
	assert.Nil(t, resp.Err)
}

func TestHttpError(t *testing.T) {
	c := NewClient(serverUrl, false, 10000)
	reqBody := ""
	resp := c.SendRequest("GET", "/", []string{}, &reqBody)
	assert.Nil(t, resp.Err)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestConnectionError(t *testing.T) {
	c := NewClient("http://localhost:9999", false, 10000)
	reqBody := ""
	resp := c.SendRequest("GET", "/potato", []string{}, &reqBody)
	assert.NotNil(t, resp.Err)
}

func setup() {
	pathResponseHandlerFunc := func(rw http.ResponseWriter, r *http.Request) {
		if want, have := "/path", r.URL.Path; want != have {
			rw.WriteHeader(404)
		}
	}
	pathHandler := fixture.PathResponseHandler{Path: WorkingPath, PathHandlerFunc: pathResponseHandlerFunc}
	var mockServerPort int
	mockServer, mockServerPort = fixture.StartHttpTargetTestServer([]fixture.PathResponseHandler{pathHandler})

	serverUrl = "http://localhost:" + fmt.Sprint(mockServerPort)
}

func teardown() {
	mockServer.Shutdown(context.Background())
}
