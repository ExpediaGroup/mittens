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

// Utils for probe files.

package probe

import (
	"io/ioutil"
	"log"
)

// WriteFile writes sample content to a file. This file can be used as a liveness/readiness check e.g. in Kubernetes.
func WriteFile(file string) {
	log.Printf("Writing file: %s", file)

	fileBytes := []byte("foo bar")

	if err := ioutil.WriteFile(file, fileBytes, 0644); err != nil {
		log.Printf("Writing to file failed with error: %v", err)
		return
	}
	log.Printf("Wrote file: %s", file)
}
