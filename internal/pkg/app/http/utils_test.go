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
	"net/http"
	"os"
	"regexp"
	"testing"

	"mittens/internal/pkg/app/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttp_FlagToHttpRequest(t *testing.T) {
	requestFlag := `post:/db:{"db": "true"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)
	assert.Equal(t, `{"db": "true"}`, *request.Body)
}

func TestBodyFromFile(t *testing.T) {
	file := internal.CreateTempFile(`{"foo": "bar"}`)

	// clean up the file at the end
	defer os.Remove(file)

	requestFlag := `post:/db:file:` + file
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)
	assert.Equal(t, `{"foo": "bar"}`, *request.Body)
}

func TestHttp_FlagWithoutBodyToHttpRequest(t *testing.T) {
	requestFlag := `get:ping`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "ping", request.Path)
	assert.Nil(t, request.Body)
}

func TestHttp_FlagWithInvalidMethodToHttpRequest(t *testing.T) {
	requestFlag := `hmm:/ping:all=true`
	_, err := ToHTTPRequest(requestFlag)
	require.Error(t, err)
}

func TestHttp_TimestampInterpolation(t *testing.T) {
	requestFlag := `post:/path_{$currentTimestamp}:{"body": "{$currentTimestamp}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var numbersRegex = regexp.MustCompile("\\d+")
	matchPath := numbersRegex.MatchString(request.Path)
	matchBody := numbersRegex.MatchString(*request.Body)

	assert.True(t, matchPath)
	assert.True(t, matchBody)
	assert.Equal(t, len(request.Path), 19)  //  "path_ + 13 numbers for timestamp
	assert.Equal(t, len(*request.Body), 25) // { "body": 13 numbers for timestamp
}

func TestHttp_Interpolation(t *testing.T) {
	requestFlag := `post:/path_{$range|min=1,max=2}_{$random|foo,bar}:{"body": "{$random|foo,bar} {$range|min=1,max=2}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var pathRegex = regexp.MustCompile(`/path_\d_(foo|bar)`)
	matchPath := pathRegex.MatchString(request.Path)

	var bodyRegex = regexp.MustCompile("{\"body\": \"(foo|bar) \\d\"}")
	matchBody := bodyRegex.MatchString(*request.Body)

	assert.True(t, matchPath)
	assert.True(t, matchBody)
}
