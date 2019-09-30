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

// Performs operations on files.

package routine

import (
	"io/ioutil"
	"log"
)

// Writes dummy content to a file. This file can be used as a liveness/readiness check in Kubernetes.
func writeFile(file string) {
	log.Printf("writing file %s\n", file)

	fileBytes := []byte("foo bar")
	err := ioutil.WriteFile(file, fileBytes, 0644)

	if err = ioutil.WriteFile(file, fileBytes, 0644); err != nil {
		log.Printf("writing to file failed with error %v\n", err)
		return
	}
	log.Printf("wrote file %s\n", file)
}
