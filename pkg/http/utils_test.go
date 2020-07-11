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
	"net/http"
	"regexp"
	"testing"
	"time"

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

func TestHttp_FlagWithoutBodyToHttpRequest(t *testing.T) {
	requestFlag := `get:ping`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "ping", request.Path)
	assert.Nil(t, request.Body)
}

func TestHttp_DateInterpolation(t *testing.T) {
	requestFlag := `post:/db_{$currentDate}:{"date": "{$currentDate|days+5,months+2,years-1}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	dateToday := time.Now().Format("2006-01-02")                        // today + 5
	dateWithOffset := time.Now().AddDate(-1, 2, 5).Format("2006-01-02") // today -1 year, +2 months, +5 days
	assert.Equal(t, "/db_"+dateToday, request.Path)
	assert.Equal(t, fmt.Sprintf(`{"date": "%s"}`, dateWithOffset), *request.Body)
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

func TestHttp_MultipleInterpolation(t *testing.T) {
	requestFlag := `post:/path_{$range|min=1,max=2}_{$random|foo,bar}:{"body": "{$random|foo,bar} {$range|min=1,max=2}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var pathRegex = regexp.MustCompile("/path_\\d_(foo|bar)")
	matchPath := pathRegex.MatchString(request.Path)

	var bodyRegex = regexp.MustCompile("{\"body\": \"(foo|bar) \\d\"}")
	matchBody := bodyRegex.MatchString(*request.Body)

	assert.True(t, matchPath)
	assert.True(t, matchBody)
}

func TestHttp_RangeInterpolation(t *testing.T) {
	requestFlag := `post:/path_{$range|min=1,max=2}:{"body": "{$range|min=1,max=2}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var pathRegex = regexp.MustCompile("/path_\\d")
	matchPath := pathRegex.MatchString(request.Path)

	var bodyRegex = regexp.MustCompile("{\"body\": \"\\d\"}")
	matchBody := bodyRegex.MatchString(*request.Body)

	assert.True(t, matchPath)
	assert.True(t, matchBody)
}

func TestHttp_InvalidRangeInterpolation(t *testing.T) {
	requestFlag := `post:/path_{$range|min=2,max=1}:{"body": "{$range|min=2,max=1}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	// will not action on invalid ranges
	assert.Equal(t, request.Path, "/path_{$range|min=2,max=1}")
	assert.Equal(t, *request.Body, "{\"body\": \"{$range|min=2,max=1}\"}")
}

func TestHttp_RandomElementInterpolation(t *testing.T) {
	requestFlag := `post:/path_{$random|fo-o,b_ar}:{"body": "{$random|fo-o,b_ar}"}`
	request, err := ToHTTPRequest(requestFlag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var elementsRegex = regexp.MustCompile("(fo-o|b_ar)")
	matchPath := elementsRegex.MatchString(request.Path)
	matchBody := elementsRegex.MatchString(*request.Body)

	assert.True(t, matchPath)
	assert.True(t, matchBody)
}
