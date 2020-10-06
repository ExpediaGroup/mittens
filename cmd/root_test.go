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

// +build integration

package cmd

import (
	grpcurl_testing "github.com/fullstorydev/grpcurl/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/reflection"
	"log"
	"mittens/pkg/safe"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

// this is a hack since Go doesn't support setup/tearDown
// we use sub-tests so that target servers only start once
func TestAll(t *testing.T) {
	shutdownHttp := StartHttpTargetTestServer(t)
	shutdownGrpc := StartGrpcTargetTestServer(t)
	defer shutdownHttp()
	defer shutdownGrpc()

	result := t.Run("TestGrpcAndHttp", TestGrpcAndHttp)
	result = result && t.Run("TestWarmupSidecarWithFileProbe", TestWarmupSidecarWithFileProbe)
	result = result && t.Run("TestWarmupFailReadinessIfTargetIsNeverReady", TestWarmupFailReadinessIfTargetIsNeverReady)
	result = result && t.Run("TestWarmupFailReadinessIfNoRequestsAreSentToTarget", TestWarmupFailReadinessIfNoRequestsAreSentToTarget)
	result = result && t.Run("TestShouldBeReadyRegardlessIfWarmupRan", TestShouldBeReadyRegardlessIfWarmupRan)
	result = result && t.Run("TestShouldBeReadyRegardlessIfHasPanicked", TestShouldBeReadyRegardlessIfHasPanicked)
	os.Exit(bool2int(!result))
}

func TestShouldBeReadyRegardlessIfWarmupRan(t *testing.T) {
	deleteFile("alive")
	deleteFile("ready")

	os.Args = []string{"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/non-existent",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=5"}

	CreateConfig()
	RunCmdRoot()

	assert.Equal(t, true, opts.FileProbe.Enabled)
	assert.ElementsMatch(t, opts.HTTP.Requests, []string{"get:/non-existent"})
	assert.Equal(t, 2, opts.Concurrency)
	assert.Equal(t, true, opts.ExitAfterWarmup)
	assert.Equal(t, "/health", opts.Target.ReadinessHTTPPath)
	assert.Equal(t, 5, opts.MaxDurationSeconds)

	readyFileExists, err := fileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

func TestShouldBeReadyRegardlessIfHasPanicked(t *testing.T) {
	deleteFile("ready")

	// we trigger a panic scenario by using a non-existent gRPC readiness probe
	os.Args = []string{"mittens",
		"-file-probe-enabled=true",
		"-exit-after-warmup=true",
		"-fail-readiness=false",
		"-target-readiness-protocol=grpc",
		"-target-grpc-port=50051",
		"-target-readiness-grpc-method=non.existent/NonExistent",
		"-target-insecure=true",
		"-max-duration-seconds=5"}

	CreateConfig()
	RunCmdRoot()

	assert.True(t, safe.HasPanicked())
	readyFileExists, err := fileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

func TestWarmupSidecarWithFileProbe(t *testing.T) {
	deleteFile("alive")
	deleteFile("ready")

	os.Args = []string{"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/delay",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=5"}

	CreateConfig()
	RunCmdRoot()

	assert.Equal(t, true, opts.FileProbe.Enabled)
	assert.ElementsMatch(t, opts.HTTP.Requests, []string{"get:/delay"})
	assert.Equal(t, 2, opts.Concurrency)
	assert.Equal(t, true, opts.ExitAfterWarmup)
	assert.Equal(t, "/health", opts.Target.ReadinessHTTPPath)
	assert.Equal(t, 5, opts.MaxDurationSeconds)

	readyFileExists, err := fileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

func TestWarmupFailReadinessIfTargetIsNeverReady(t *testing.T) {
	deleteFile("alive")
	deleteFile("ready")

	// we simulate a failure in the target by setting the readiness path to a non existent one so that
	// the target never becomes ready and the warmup does not run
	os.Args = []string{"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/delay",
		"-target-readiness-port=8080",
		"-target-readiness-http-path=/non-existent",
		"-max-duration-seconds=5",
		"-exit-after-warmup=true",
		"-fail-readiness=true"}

	CreateConfig()
	RunCmdRoot()

	assert.Equal(t, true, opts.FileProbe.Enabled)
	assert.ElementsMatch(t, opts.HTTP.Requests, []string{"get:/delay"})
	assert.Equal(t, true, opts.ExitAfterWarmup)
	assert.Equal(t, "/non-existent", opts.Target.ReadinessHTTPPath)
	assert.Equal(t, 5, opts.MaxDurationSeconds)
	assert.Equal(t, true, opts.FailReadiness)

	readyFileExists, err := fileExists("ready")
	require.NoError(t, err)
	assert.False(t, readyFileExists)
}

func TestWarmupFailReadinessIfNoRequestsAreSentToTarget(t *testing.T) {
	deleteFile("alive")
	deleteFile("ready")

	// we simulate a failure by using a port that doesnt exist (9999)
	os.Args = []string{"mittens",
		"-file-probe-enabled=true",
		"-http-requests=get:/delay",
		"-target-http-port=9999",
		"-target-readiness-port=8080",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=5",
		"-exit-after-warmup=true",
		"-fail-readiness=true"}

	CreateConfig()
	RunCmdRoot()

	assert.Equal(t, true, opts.FileProbe.Enabled)
	assert.ElementsMatch(t, opts.HTTP.Requests, []string{"get:/delay"})
	assert.Equal(t, true, opts.ExitAfterWarmup)
	assert.Equal(t, "/health", opts.Target.ReadinessHTTPPath)
	assert.Equal(t, 5, opts.MaxDurationSeconds)
	assert.Equal(t, true, opts.FailReadiness)

	readyFileExists, err := fileExists("ready")
	require.NoError(t, err)
	assert.False(t, readyFileExists)
}

func TestGrpcAndHttp(t *testing.T) {
	deleteFile("alive")
	deleteFile("ready")

	os.Args = []string{"mittens",
		"-file-probe-enabled=true",
		"-target-grpc-port=50051",
		"-http-requests=get:/delay",
		"-grpc-requests=grpc.testing.TestService/EmptyCall",
		"-grpc-requests=grpc.testing.TestService/UnaryCall:{\"payload\":{\"body\":\"abcdefghijklmnopqrstuvwxyz01\"}}",
		"-target-insecure=true",
		"-concurrency=2",
		"-exit-after-warmup=true",
		"-target-readiness-http-path=/health",
		"-max-duration-seconds=5"}

	CreateConfig()
	RunCmdRoot()

	assert.Equal(t, true, opts.FileProbe.Enabled)
	assert.ElementsMatch(t, opts.HTTP.Requests, []string{"get:/delay"})
	assert.ElementsMatch(t, opts.Grpc.Requests, []string{"grpc.testing.TestService/EmptyCall", "grpc.testing.TestService/UnaryCall:{\"payload\":{\"body\":\"abcdefghijklmnopqrstuvwxyz01\"}}"})

	assert.Equal(t, 2, opts.Concurrency)
	assert.Equal(t, true, opts.ExitAfterWarmup)
	assert.Equal(t, "/health", opts.Target.ReadinessHTTPPath)
	assert.Equal(t, 5, opts.MaxDurationSeconds)

	readyFileExists, err := fileExists("ready")
	require.NoError(t, err)
	assert.True(t, readyFileExists)
}

// StartGrpcTargetTestServer starts a gRPC server in port 50051
// It uses the test.proto from grpc-testing: https://github.com/grpc/grpc-go/blob/40a879c23a0dc77234d17e0699d074d5fd151bd0/test/grpc_testing/test.proto
func StartGrpcTargetTestServer(t *testing.T) (shutdown func()) {
	svr := grpc.NewServer()

	grpc_testing.RegisterTestServiceServer(svr, grpcurl_testing.TestServer{})
	reflection.Register(svr)
	l, _ := net.Listen("tcp", ":50051")

	var serverErr error
	go func() {
		serverErr = svr.Serve(l)
	}()

	time.Sleep(100 * time.Millisecond)
	require.NoError(t, serverErr)

	return shutdown
}

func StartHttpTargetTestServer(t *testing.T) (shutdown func()) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/delay", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.WriteHeader(http.StatusNoContent)
	})

	server := &http.Server{Addr: ":8080"}

	var serverErr error
	go func() {
		serverErr = server.ListenAndServe()
	}()

	// wait for server to star up
	time.Sleep(100 * time.Millisecond)
	require.NoError(t, serverErr)

	return shutdown
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func deleteFile(path string) {
	var err = os.Remove(path)
	if err != nil {
		log.Printf("File not deleted")
	}
}

func fileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
