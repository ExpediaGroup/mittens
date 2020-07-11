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

import "net/http"

// Handler stores the liveness and readiness values for the probes
type Handler struct {
	alive bool
	ready bool
}

// isAlive sets the liveness probe
func (h *Handler) isAlive(alive bool) {
	h.alive = alive
}

// IsReady sets the readiness probe
func (h *Handler) IsReady(ready bool) {
	h.ready = ready
}

// aliveHandler returns 200 if the probe server is alive and 404 if it is not
func (h *Handler) aliveHandler() func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if h.alive {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

// aliveHandler returns 200 if the probe server is ready and 404 if it is not
func (h *Handler) readyHandler() func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if h.ready {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
