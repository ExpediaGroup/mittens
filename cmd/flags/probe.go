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

type Probe struct {
	Port          int
	LivenessPath  string
	ReadinessPath string
}

func (p *Probe) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *Probe) InitFlags(cmd *cobra.Command) {

	cmd.Flags().IntVar(
		&p.Port,
		"probe-port",
		8000,
		"Warm up sidecar port for liveness and readiness probe",
	)
	cmd.Flags().StringVar(
		&p.LivenessPath,
		"probe-liveness-path",
		"/alive",
		"Warm up sidecar liveness probe path",
	)
	cmd.Flags().StringVar(
		&p.ReadinessPath,
		"probe-readiness-path",
		"/ready",
		"Warm up sidecar readiness probe path",
	)
}