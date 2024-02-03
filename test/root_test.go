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

package test

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"mittens/cmd"
	"mittens/fixture"
	"mittens/internal/pkg/probe"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var mockHttpServerPort int
var mockHttpServer *http.Server
var mockGrpcServer *grpc.Server
var httpInvocations = 0
var decompressedBody = ""

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func TestShouldBeReadyRegardlessIfWarmupRan(t *testing.T) {
	t.Cleanup(func() {
		cleanup()
	})

	os.Args = []string{
		"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/non-existent",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=2",
	}

	cmd.CreateConfig()
	cmd.RunCmdRoot()

	assert.Equal(t, httpInvocations, 0, "Assert that no calls were made to the http service")

	readyFileExists, err := probe.FileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

func TestWarmupFailReadinessIfTargetIsNeverReady(t *testing.T) {
	t.Cleanup(func() {
		cleanup()
	})

	// we simulate a failure in the target by setting the readiness path to a non existent one so that
	// the target never becomes ready and the warmup does not run
	os.Args = []string{
		"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/hello-world",
		fmt.Sprintf("-target-http-port=%d", mockHttpServerPort),
		fmt.Sprintf("-target-readiness-port=%d", mockHttpServerPort),
		"-target-readiness-http-path=/non-existent",
		"-max-duration-seconds=2",
		"-exit-after-warmup=true",
		"-fail-readiness=true",
	}

	cmd.CreateConfig()
	cmd.RunCmdRoot()

	assert.Equal(t, httpInvocations, 0, "Assert that no calls were made to the http service")

	readyFileExists, err := probe.FileExists("ready")
	require.NoError(t, err)
	assert.False(t, readyFileExists)
}

func TestWarmupFailReadinessIfNoRequestsAreSentToTarget(t *testing.T) {
	t.Cleanup(func() {
		cleanup()
	})

	// we simulate a failure by using a port that doesnt exist (9999)
	os.Args = []string{
		"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/hello-world",
		"-target-http-port=9999",
		fmt.Sprintf("-target-readiness-port=%d", mockHttpServerPort),
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=2",
		"-exit-after-warmup=true",
		"-fail-readiness=true",
	}

	cmd.CreateConfig()
	cmd.RunCmdRoot()

	assert.Equal(t, httpInvocations, 0, "Assert that no calls were made to the http service")

	readyFileExists, err := probe.FileExists("ready")
	require.NoError(t, err)
	assert.False(t, readyFileExists)
}

func TestHttp(t *testing.T) {
	t.Cleanup(func() {
		cleanup()
	})

	os.Args = []string{
		"mittens",
		"-file-probe-enabled=true",
		// FIXME: for some reason we need to set both ports?
		fmt.Sprintf("-target-http-port=%d", mockHttpServerPort),
		fmt.Sprintf("-target-readiness-port=%d", mockHttpServerPort),
		"-http-requests=get:/hello-world",
		"-target-insecure=true",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=2",
		"-concurrency-target-seconds=1",
	}

	cmd.CreateConfig()
	cmd.RunCmdRoot()

	assert.Greater(t, httpInvocations, 1, "Assert that we made some calls to the http service")
	// TODO: validate grpc invocations

	readyFileExists, err := probe.FileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

func TestGrpcAndHttp(t *testing.T) {
	t.Cleanup(func() {
		cleanup()
	})

	os.Args = []string{
		"mittens",
		"-file-probe-enabled=true",
		"-target-grpc-port=50051",
		// FIXME: for some reason we need to set both ports?
		fmt.Sprintf("-target-http-port=%d", mockHttpServerPort),
		fmt.Sprintf("-target-readiness-port=%d", mockHttpServerPort),
		"-http-requests=get:/hello-world",
		"-grpc-requests=grpc.testing.TestService/EmptyCall",
		"-grpc-requests=grpc.testing.TestService/UnaryCall:{\"payload\":{\"body\":\"abcdefghijklmnopqrstuvwxyz01\"}}",
		"-target-insecure=true",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=2",
		"-concurrency-target-seconds=1",
	}

	cmd.CreateConfig()
	cmd.RunCmdRoot()

	assert.Greater(t, httpInvocations, 1, "Assert that we made some calls to the http service")
	// TODO: validate grpc invocations

	readyFileExists, err := probe.FileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

func TestCompressWithGZip(t *testing.T) {
	t.Cleanup(func() {
		cleanup()
	})

	var requestBody = "{\"payload\":{\"body\":\"abcdefghijklmnopqrstuvwxyz01\"}}"

	os.Args = []string{
		"mittens",
		"-file-probe-enabled=true",
		// FIXME: for some reason we need to set both ports?
		fmt.Sprintf("-target-http-port=%d", mockHttpServerPort),
		fmt.Sprintf("-target-readiness-port=%d", mockHttpServerPort),
		"-http-requests=post:/compressed:" + requestBody,
		"-http-requests-compression=gzip",
		"-target-insecure=true",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=2",
		"-concurrency-target-seconds=1",
	}

	cmd.CreateConfig()
	cmd.RunCmdRoot()

	assert.Greater(t, httpInvocations, 1, "Assert that we made some calls to the http service")
	assert.Equal(t, requestBody, decompressedBody, "Assert that server-side decompressed body is equal to client request body")
}

func setup() {
	fmt.Println("Starting up http server")
	mockHttpServer, mockHttpServerPort = fixture.StartHttpTargetTestServer([]fixture.PathResponseHandler{
		{
			Path: "/hello-world",
			PathHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				// Tiny sleep to mimic a regular http call
				time.Sleep(time.Millisecond * 10)
				// Record number of invocations made to this endpoint
				httpInvocations++
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			Path: "/compressed",
			PathHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				// Tiny sleep to mimic a regular http call
				time.Sleep(time.Millisecond * 10)
				httpInvocations++
				// Record number of invocations made to this endpoint
				if r.Header.Get("Content-Encoding") == "gzip" {
					gr, err := gzip.NewReader(r.Body)
					if err != nil {
						 w.WriteHeader(http.StatusInternalServerError)
						 return
					}
					b, err := io.ReadAll(gr)
					if err != nil {
						 w.WriteHeader(http.StatusInternalServerError)
						 return
					}
					decompressedBody = string(b)
					gr.Close()
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
		},
	})

	// FIXME: should run on a random/free port
	fmt.Println("Starting up grpc server")
	mockGrpcServer = fixture.StartGrpcTargetTestServer(50051)
}

func teardown() {
	fmt.Println("Shutting down http server")
	mockHttpServer.Shutdown(context.Background())
	fmt.Println("Shutting down grpc server")
	mockGrpcServer.GracefulStop()

	fmt.Println("All servers server stopped")
}

func cleanup() {
	httpInvocations = 0

	if fileExists, err := probe.FileExists("alive"); err == nil && fileExists {
		probe.DeleteFile("alive")
	}
	if fileExists, err := probe.FileExists("ready"); err == nil && fileExists {
		probe.DeleteFile("ready")
	}
}
