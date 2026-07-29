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

package cmd

import (
	"mittens/cmd/flags"
	"os"
	"syscall"
	"testing"
	"time"
)

// TestBlockReturnsWhenExitAfterWarmup verifies that block() is a no-op when
// -exit-after-warmup is set, so the process can exit right after warm-up.
func TestBlockReturnsWhenExitAfterWarmup(t *testing.T) {
	prev := opts
	t.Cleanup(func() { opts = prev })

	opts = &flags.Root{}
	opts.ExitAfterWarmup = true

	done := make(chan struct{})
	go func() {
		block()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("block did not return with -exit-after-warmup set")
	}
}

// TestShutdownSignalsAreTermination pins the signal set block() registers for,
// so swapping in the wrong signal is caught even though os/signal offers no way
// to introspect a channel's registration.
func TestShutdownSignalsAreTermination(t *testing.T) {
	want := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	if len(shutdownSignals) != len(want) {
		t.Fatalf("shutdownSignals = %v, want %v", shutdownSignals, want)
	}
	for i, s := range want {
		if shutdownSignals[i] != s {
			t.Errorf("shutdownSignals[%d] = %v, want %v", i, shutdownSignals[i], s)
		}
	}
}

// TestBlockUntilSignalUnblocksOnSignal verifies the keep-alive wait is
// wakeable: it returns once a termination signal is delivered. This guards
// against regressing to a bare `select {}`, which is unwakeable and trips
// Go's deadlock detector when the process goes idle (see #366).
func TestBlockUntilSignalUnblocksOnSignal(t *testing.T) {
	sig := make(chan os.Signal, 1)

	done := make(chan struct{})
	go func() {
		blockUntilSignal(sig)
		close(done)
	}()

	sig <- syscall.SIGTERM

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("blockUntilSignal did not return after receiving a signal")
	}
}
