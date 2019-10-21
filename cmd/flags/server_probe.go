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
	"fmt"
	"github.com/spf13/cobra"
)

type ServerProbe struct {
	Enabled       bool
	Port          int
	LivenessPath  string
	ReadinessPath string
}

func (p *ServerProbe) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *ServerProbe) InitFlags(cmd *cobra.Command) {

	cmd.Flags().BoolVar(
		&p.Enabled,
		"probe-server-enabled",
		false,
		"If set to true runs a web server that exposes endpoints to be used as readiness/liveness probes",
	)
	cmd.Flags().IntVar(
		&p.Port,
		"probe-server-port",
		8000,
		"Port on which probe server is running",
	)
	cmd.Flags().StringVar(
		&p.LivenessPath,
		"probe-server-liveness-path",
		"/alive",
		"Probe server endpoint used as liveness probe",
	)
	cmd.Flags().StringVar(
		&p.ReadinessPath,
		"probe-server-readiness-path",
		"/ready",
		"Probe server endpoint used as readiness probe",
	)
}
