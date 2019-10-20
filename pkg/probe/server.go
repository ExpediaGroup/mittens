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

// Http server with addition 'IsAlive' and 'IsReady' methods
type Server struct {
	*Handler
	httpServer *http.Server
}

func NewServer(port int, livenessPath, readinessPath string) *Server {

	handler := new(Handler)

	mux := http.NewServeMux()
	mux.HandleFunc(livenessPath, handler.aliveHandler())
	mux.HandleFunc(readinessPath, handler.readyHandler())

	log.Printf("probe server on %d port", port)
	log.Printf("liveness path: %s", livenessPath)
	log.Printf("readiness path: %s", readinessPath)

	return &Server{
		httpServer: newServer(port, mux),
		Handler:    handler,
	}
}

func (s *Server) ListenAndServe() error {

	s.isAlive(true)
	log.Print("starting probe server")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() {

	s.IsReady(false)
	s.isAlive(false)
	log.Print("shutting down probe server")
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("probe server shutdown: %v", err)
	}
}

func newServer(port int, handler http.Handler) *http.Server {

	return &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
