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

package probe

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_IsAliveCanBeSet(t *testing.T) {

	handler := new(Handler)
	server := httptest.NewServer(http.HandlerFunc(handler.aliveHandler()))
	defer server.Close()

	// livness probe defaults to false
	resp, err := http.DefaultClient.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	// set livness probe to true
	handler.isAlive(true)
	resp, err = http.DefaultClient.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// set livness probe again to false
	handler.isAlive(false)
	resp, err = http.DefaultClient.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestServer_IsReadyCanBeSet(t *testing.T) {

	handler := new(Handler)
	server := httptest.NewServer(http.HandlerFunc(handler.readyHandler()))
	defer server.Close()

	// readiness probe defaults to false
	resp, err := http.DefaultClient.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	// set readiness probe to true
	handler.IsReady(true)
	resp, err = http.DefaultClient.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// set readiness probe again to false
	handler.IsReady(false)
	resp, err = http.DefaultClient.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
