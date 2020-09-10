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

package grpc

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrpc_FlagToGrpcRequest(t *testing.T) {

	requestFlag := `health/ping:{"db": "true"}`
	request, err := ToGrpcRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, "health/ping", request.ServiceMethod)
	assert.Equal(t, `{"db": "true"}`, string(request.Message))
}

func TestGrpc_FlagWithoutBodyToGrpcRequest(t *testing.T) {

	requestFlag := `health/ping`
	request, err := ToGrpcRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, "health/ping", request.ServiceMethod)
	assert.Equal(t, "", string(request.Message))
}

func TestGrpc_InvalidFlagToGrpcRequest(t *testing.T) {

	requestFlag := `health:ping`
	_, err := ToGrpcRequest(requestFlag)
	require.Error(t, err)
}

func TestGrpc_Interpolation(t *testing.T) {
	requestFlag := `health/ping:{"lorem": "{$random|foo}", "ipsum":"{$random|foo}"}"}`
	request, err := ToGrpcRequest(requestFlag)
	require.NoError(t, err)

	var pathRegex = regexp.MustCompile(`{"lorem": "(foo|bar)", "ipsum":"(foo|bar)"}`)
	matchRequest := pathRegex.MatchString(request.Message)

	assert.True(t, matchRequest)
}
