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
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

// this is a hack since Go doesn't support setup/tearDown
// we use sub-tests so that target server only starts once

var _ = func() bool {
	testing.Init()
	return true
}()

func TestAll(t *testing.T) {
	shutdown := StartTargetTestServer(t)
	defer shutdown()
	result := t.Run("TestWarmupSidecarWithFileProbe", TestWarmupSidecarWithFileProbe)
	result = result && t.Run("TestWarmupSidecarWithServerProbe", TestWarmupSidecarWithServerProbe)
	os.Exit(bool2int(!result))
}

func TestWarmupSidecarWithFileProbe(t *testing.T) {
	os.Args = []string{"mittens",
		"-fileProbeEnabled=true",
		"-serverProbeEnabled=false",
		"-httpRequest=get:/delay",
		"-concurrency=4",
		"-exitAfterWarmup=true",
		"-targetReadinessPath=/",
		"-timeoutSeconds=5"}

	RunCmdRoot()
}

func TestWarmupSidecarWithServerProbe(t *testing.T) {
	os.Args = []string{"mittens",
		"-fileProbeEnabled=true",
		"-serverProbeEnabled=true",
		"-httpRequest=get:/delay",
		"-concurrency=4",
		"-exitAfterWarmup=true",
		"-targetReadinessPath=/",
		"-timeoutSeconds=5"}
	RunCmdRoot()
}

func StartTargetTestServer(t *testing.T) (shutdown func()) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		log.Print("handler /")
	})

	http.HandleFunc("/delay/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.WriteHeader(http.StatusNoContent)
		log.Print("handler /delay/")
	})

	server := &http.Server{Addr: ":8080"}
	shutdown = func() {
		err := server.Shutdown(context.Background())
		assert.NoError(t, err)
	}

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
