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
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server represents an HTTP server with liveness and readiness probes.
type Server struct {
	*Handler
	httpServer *http.Server
}

// NewServer creates a Server instance with liveness and readiness handlers.
func NewServer(port int, livenessPath, readinessPath string) *Server {

	handler := new(Handler)

	mux := http.NewServeMux()
	mux.HandleFunc(livenessPath, handler.aliveHandler())
	mux.HandleFunc(readinessPath, handler.readyHandler())

	log.Printf("Probe server on %d port", port)
	log.Printf("Liveness path: %s", livenessPath)
	log.Printf("Readiness path: %s", readinessPath)

	return &Server{
		httpServer: newServer(port, mux),
		Handler:    handler,
	}
}

// ListenAndServe starts the probe server and enables the liveness probe.
func (s *Server) ListenAndServe() error {
	s.isAlive(true)
	log.Print("Starting probe server")
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the probe server and disables the liveness and readiness probes.
func (s *Server) Shutdown() {
	s.IsReady(false)
	s.isAlive(false)
	log.Print("Shutting down probe server")
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Probe server shutdown: %v", err)
	}
}

// newServer returns an HTTP server.
func newServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
