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

package flags

import (
	"flag"
	"fmt"
)

type FailReadiness struct {
	Enabled      bool `json:"fail-readiness-enabled"`
	NoRequests   bool `json:"fail-readiness-no-requests"`
	ClientErrors bool `json:"fail-readiness-client-errors"`
}

func (p *FailReadiness) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *FailReadiness) InitFlags() {
	flag.BoolVar(&p.Enabled, "fail-readiness-enabled", false, "If set to true readiness will fail if there were issues with the warmup")
	flag.BoolVar(&p.NoRequests, "fail-readiness-no-requests", true, "Readiness will fail if no requests were sent")
	flag.BoolVar(&p.ClientErrors, "fail-readiness-client-errors", false, "Readiness will fail in case of client errors")
}
