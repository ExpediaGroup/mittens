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
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestSuccess(t *testing.T) {
	path := "/path"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if want, have := path, r.URL.Path; want != have {
			t.Errorf("unexpected path, want: %q, have %q", want, have)
		}
	}))
	defer server.Close()

	c := NewClient(server.URL, false)
	reqBody := ""
	resp := c.Request("GET", path, map[string]string{}, &reqBody)
	assert.True(t, resp.RequestSent)
	assert.Nil(t, resp.Err)
}

func TestClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(400)
	}))
	defer server.Close()

	c := NewClient(server.URL, false)
	reqBody := ""
	resp := c.Request("GET", "/", map[string]string{}, &reqBody)
	assert.True(t, resp.RequestSent)
	assert.NotNil(t, resp.Err)
}

func TestServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(500)
	}))
	defer server.Close()

	c := NewClient(server.URL, false)
	reqBody := ""
	resp := c.Request("GET", "/", map[string]string{}, &reqBody)
	assert.True(t, resp.RequestSent)
	assert.NotNil(t, resp.Err)
}

func TestRequestNotSent(t *testing.T) {
	c := NewClient("/", false)
	reqBody := ""
	resp := c.Request("GET", "/", map[string]string{}, &reqBody)
	assert.False(t, resp.RequestSent)
	assert.NotNil(t, resp.Err)
}
