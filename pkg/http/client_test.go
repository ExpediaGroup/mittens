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
	"fmt"
	"github.com/stretchr/testify/assert"
	"mittens/fixture"
	"net/http"
	"testing"
)

var mockServer *http.Server

const WorkingPath = "/path"
const Port int = 8090

var ServerUrl = "http://localhost:" + fmt.Sprint(Port)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func TestRequestSuccess(t *testing.T) {
	c := NewClient(ServerUrl, false)
	reqBody := ""
	resp := c.SendRequest("GET", WorkingPath, map[string]string{}, &reqBody)
	assert.Nil(t, resp.Err)
}

func TestHttpError(t *testing.T) {
	c := NewClient(ServerUrl, false)
	reqBody := ""
	resp := c.SendRequest("GET", "/", map[string]string{}, &reqBody)
	assert.Nil(t, resp.Err)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestConnectionError(t *testing.T) {
	c := NewClient("http://localhost:9999", false)
	reqBody := ""
	resp := c.SendRequest("GET", "/potato", map[string]string{}, &reqBody)
	assert.NotNil(t, resp.Err)
}

func setup() {
	pathReposnseHandlerFunc := func(rw http.ResponseWriter, r *http.Request) {
		if want, have := "/path", r.URL.Path; want != have {
			rw.WriteHeader(404)
		}
	}
	pathHandler := fixture.PathResponseHandler{WorkingPath, pathReposnseHandlerFunc}
	mockServer = fixture.StartHttpTargetTestServer(Port, []fixture.PathResponseHandler{pathHandler}, false)
}

func teardown() {
	mockServer.Close()
}
