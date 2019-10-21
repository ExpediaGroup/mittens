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
	"net/http"
	"testing"
)

func TestHttp_ToHttpRequests(t *testing.T) {

	requestFlags := []string{
		"get:/health",
		"get:/ping",
	}

	requests, err := toHttpRequests(requestFlags)
	require.NoError(t, err)

	require.Equal(t, 2, len(requests))
	assert.Equal(t, "/health", requests[0].Path)
	assert.Equal(t, "/ping", requests[1].Path)
}

func TestHttp_FlagToHttpRequest(t *testing.T) {

	requestFlag := `post:/db:{"db": "true"}`
	request, err := toHttpRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)
	assert.Equal(t, `{"db": "true"}`, string(request.Body))
}

func TestHttp_FlagWithoutBodyToHttpRequest(t *testing.T) {

	requestFlag := `get:ping`
	request, err := toHttpRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "ping", request.Path)
	assert.Nil(t, request.Body)
}

func TestHttp_FlagWithInvlidMethodToHttpRequest(t *testing.T) {

	requestFlag := `hmm:/ping:all=true`
	_, err := toGrpcRequest(requestFlag)
	require.Error(t, err)
}
