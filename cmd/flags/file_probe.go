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

type FileProbe struct {
	Enabled       bool `json:"enabled"`
	LivenessPath  string `json:"livenessPath"`
	ReadinessPath string `json:"readinessPath"`
}

func (p *FileProbe) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *FileProbe) InitFlags() {
	//flag.BoolVar(&p.Enabled,"probe-file-enabled",true,"If set to true writes files to be used as readiness/liveness probes")
	flag.StringVar(&p.LivenessPath,"fileProbeLivenessPath","alive","File to be used for liveness probe")
	flag.StringVar(&p.ReadinessPath,"fileProbeReadinessPath","ready","File to be used for readiness probe")
}
