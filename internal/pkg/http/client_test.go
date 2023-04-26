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
const AuthTokenPath = "/mittens/token"

var serverUrl string

func TestMain(m *testing.M) {
	setup()

	m.Run()
	teardown()
}

func TestRequestSuccess(t *testing.T) {
	c := NewClient(serverUrl, false)
	reqBody := ""
	resp := c.SendRequest("GET", WorkingPath, []string{}, &reqBody)
	assert.Nil(t, resp.Err)
}

func TestHttpError(t *testing.T) {
	c := NewClient(serverUrl, false)
	reqBody := ""
	resp := c.SendRequest("GET", "/", []string{}, &reqBody)
	assert.Nil(t, resp.Err)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestConnectionError(t *testing.T) {
	c := NewClient("http://localhost:9999", false)
	reqBody := ""
	resp := c.SendRequest("GET", "/potato", []string{}, &reqBody)
	assert.NotNil(t, resp.Err)
}

func TestAuthRequestSuccess(t *testing.T) {
	c := NewClient(serverUrl, false)
	resp, err := c.SendAuthTokenRequest(AuthTokenPath)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestAuthRequestHttpError(t *testing.T) {
	c := NewClient(serverUrl, false)
	resp, err := c.SendAuthTokenRequest("/wrong/auth/token/path")
	assert.NotNil(t, err)
	assert.Empty(t, resp)
}

func TestAuthTokenConnectionError(t *testing.T) {
	c := NewClient("http://localhost:9999", false)
	resp, err := c.SendAuthTokenRequest("/potato")
	assert.NotNil(t, err)
	assert.Empty(t, resp)
}

func setup() {
	pathResponseHandlerFunc := func(rw http.ResponseWriter, r *http.Request) {
		if want, have := "/path", r.URL.Path; want != have {
			rw.WriteHeader(404)
		}
	}

	authResponseHandlerFunc := func(rw http.ResponseWriter, r *http.Request) {
		if want, have := "/mittens/token", r.URL.Path; want != have {
			rw.WriteHeader(404)
		}
	}

	pathHandler := fixture.PathResponseHandler{Path: WorkingPath, PathHandlerFunc: pathResponseHandlerFunc}
	authHandler := fixture.PathResponseHandler{Path: AuthTokenPath, PathHandlerFunc: authResponseHandlerFunc}
	var mockServerPort int
	mockServer, mockServerPort = fixture.StartHttpTargetTestServer([]fixture.PathResponseHandler{pathHandler, authHandler})

	serverUrl = "http://localhost:" + fmt.Sprint(mockServerPort)
}

func teardown() {
	mockServer.Shutdown(context.Background())
}
