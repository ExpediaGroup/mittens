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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mittens/internal/pkg/http"
	"testing"
)

func TestHttp_ToHttpRequests(t *testing.T) {
	requestFlags := []string{
		"get:/health",
		"get:/ping",
	}

	requests, err := toHTTPRequests(requestFlags)
	require.NoError(t, err)

	require.Equal(t, 2, len(requests))
	assert.Equal(t, "/health", requests[0].Path)
	assert.Equal(t, "/ping", requests[1].Path)
}

func TestHttp_ToHttpRequestsInvalidFormat(t *testing.T) {

	requestFlags := []string{
		"get/health",
	}

	requests, err := toHTTPRequests(requestFlags)

	var expected []http.Request
	require.Error(t, err)
	require.Equal(t, "invalid request flag: get/health, expected format <http-method>:<path>[:body]", err.Error())
	require.Equal(t, expected, requests)
}

func TestHttp_ToHttpRequestsInvalidMethod(t *testing.T) {

	requestFlags := []string{
		"invalidMethod:/health",
	}

	requests, err := toHTTPRequests(requestFlags)

	var expected []http.Request
	require.Error(t, err)
	require.Equal(t, "invalid request flag: invalidMethod:/health, method INVALIDMETHOD is not supported", err.Error())
	require.Equal(t, expected, requests)
}

func TestHttp_ToHttpRequestsInvalidBody(t *testing.T) {

	requestFlags := []string{
		"get:/test:file:test",
	}

	requests, err := toHTTPRequests(requestFlags)

	var expected []http.Request
	require.Error(t, err)
	require.Equal(t, "unable to parse body for request: file:test", err.Error())
	require.Equal(t, expected, requests)
}
