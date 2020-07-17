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
	"log"
	"strings"
)

// toHeaders converts the headers from the format these are passed by the user to a map.
func toHeaders(headersFlag []string) map[string]string {

	headers := make(map[string]string)
	for _, h := range headersFlag {
		kv := strings.SplitN(h, ":", 2)
		if len(kv) == 1 {
			log.Printf("cannot find ':' separator in supplied header %s", h)
			headers[strings.TrimSpace(kv[0])] = ""
			continue
		}
		headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return headers
}
