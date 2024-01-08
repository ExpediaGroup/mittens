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
	"bytes"
	"net/http"
	"os"
	"regexp"
	"testing"

	"mittens/internal/pkg/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttp_FlagToHttpRequest(t *testing.T) {
	requestFlag := `post:/db:{"db": "true"}`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_NONE)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)
	body := new(bytes.Buffer)
	body.ReadFrom(request.Body)
	assert.Equal(t, `{"db": "true"}`, body.String())
}

func TestHttp_CompressGzip(t *testing.T) {
	requestFlag := `post:/db:{"db": "true"}`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_GZIP)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)

	assert.Equal(t, map[string]string{"Content-Encoding": "gzip"}, request.Headers)

	body := new(bytes.Buffer)
	body.ReadFrom(request.Body)
	expected := &bytes.Buffer{}
	expected.Write([]byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xaa, 0x56, 0x4a, 0x49, 0x52, 0xb2, 0x52, 0x50, 0x2a, 0x29, 0x2a, 0x4d, 0x55, 0xaa, 0x5, 0x4, 0x0, 0x0, 0xff, 0xff, 0xa1, 0x4a, 0x9b, 0x5d, 0xe, 0x0, 0x0, 0x0})
	assert.Equal(t, expected, body)
}

func TestHttp_CompressBrotli(t *testing.T) {
	requestFlag := `post:/db:{"db": "true"}`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_BROTLI)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)

	assert.Equal(t, map[string]string{"Content-Encoding": "br"}, request.Headers)

	body := new(bytes.Buffer)
	body.ReadFrom(request.Body)
	expected := &bytes.Buffer{}
	expected.Write([]byte{0x8b, 0x6, 0x80, 0x7b, 0x22, 0x64, 0x62, 0x22, 0x3a, 0x20, 0x22, 0x74, 0x72, 0x75, 0x65, 0x22, 0x7d, 0x3})
	assert.Equal(t, expected, body)
}

func TestHttp_CompressDeflate(t *testing.T) {
	requestFlag := `post:/db:{"db": "true"}`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_DEFLATE)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)
	assert.Equal(t, map[string]string{"Content-Encoding": "deflate"}, request.Headers)
	body := new(bytes.Buffer)
	body.ReadFrom(request.Body)
	expected := &bytes.Buffer{}
	expected.Write([]byte{0xaa, 0x56, 0x4a, 0x49, 0x52, 0xb2, 0x52, 0x50, 0x2a, 0x29, 0x2a, 0x4d, 0x55, 0xaa, 0x5, 0x4, 0x0, 0x0, 0xff, 0xff})
	assert.Equal(t, expected, body)
}

func TestBodyFromFile(t *testing.T) {
	file := internal.CreateTempFile(`{"foo": "bar"}`)

	// clean up the file at the end
	defer os.Remove(file)

	requestFlag := `post:/db:file:` + file
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_NONE)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/db", request.Path)
	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	assert.Equal(t, `{"foo": "bar"}`, buf.String())
}

func TestHttp_FlagWithoutBodyToHttpRequest(t *testing.T) {
	requestFlag := `get:ping`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_NONE)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "ping", request.Path)
	assert.Nil(t, request.Body)
}

func TestHttp_FlagWithInvalidMethodToHttpRequest(t *testing.T) {
	requestFlag := `hmm:/ping:all=true`
	_, err := ToHTTPRequest(requestFlag, COMPRESSION_NONE)
	require.Error(t, err)
}

func TestHttp_TimestampInterpolation(t *testing.T) {
	requestFlag := `post:/path_{$currentTimestamp}:{"body": "{$currentTimestamp}"}`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_NONE)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var numbersRegex = regexp.MustCompile("\\d+")
	matchPath := numbersRegex.MatchString(request.Path)
	body := new(bytes.Buffer)
	body.ReadFrom(request.Body)
	matchBody := numbersRegex.MatchString(body.String())

	assert.True(t, matchPath)
	assert.True(t, matchBody)
	assert.Equal(t, len(request.Path), 19)  //  "path_ + 13 numbers for timestamp
	assert.Equal(t, len(body.String()), 25) // { "body": 13 numbers for timestamp
}

func TestHttp_Interpolation(t *testing.T) {
	requestFlag := `post:/path_{$range|min=1,max=2}_{$random|foo,bar}:{"body": "{$random|foo,bar} {$range|min=1,max=2}"}`
	request, err := ToHTTPRequest(requestFlag, COMPRESSION_NONE)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, request.Method)

	var pathRegex = regexp.MustCompile(`/path_\d_(foo|bar)`)
	matchPath := pathRegex.MatchString(request.Path)

	var bodyRegex = regexp.MustCompile("{\"body\": \"(foo|bar) \\d\"}")
	body := new(bytes.Buffer)
	body.ReadFrom(request.Body)
	matchBody := bodyRegex.MatchString(body.String())

	assert.True(t, matchPath)
	assert.True(t, matchBody)
}
