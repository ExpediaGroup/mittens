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
	"testing"
)

func TestGrpc_ToGrpcRequests(t *testing.T) {

	requestFlags := []string{
		"svc1/ping",
		"svc2/ping",
	}

	requests, err := toGrpcRequests(requestFlags)
	require.NoError(t, err)

	require.Equal(t, 2, len(requests))
	assert.Equal(t, "svc1/ping", requests[0].ServiceMethod)
	assert.Equal(t, "svc2/ping", requests[1].ServiceMethod)
}

func TestGrpc_FlagToGrpcRequest(t *testing.T) {

	requestFlag := `health/ping:{"db": "true"}`
	request, err := toGrpcRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, "health/ping", request.ServiceMethod)
	assert.Equal(t, `{"db": "true"}`, string(request.Message))
}

func TestGrpc_FlagWithoutBodyToGrpcRequest(t *testing.T) {

	requestFlag := `health/ping`
	request, err := toGrpcRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, "health/ping", request.ServiceMethod)
	assert.Nil(t, request.Message)
}

func TestGrpc_InvalidFlagToGrpcRequest(t *testing.T) {

	requestFlag := `health:ping`
	_, err := toGrpcRequest(requestFlag)
	require.Error(t, err)
}