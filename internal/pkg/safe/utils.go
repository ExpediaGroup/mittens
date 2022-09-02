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

package safe

import (
	"log"
)

// Do wraps a function with recover logic to catch unexpected panics.
func Do(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Unexpected panic was caught:", err)
		}
	}()
	f()
}

// DoAndReturn wraps a function with recover logic to catch unexpected panics.
// It returns the result of the function if no panic occurred, or the fallback result otherwise.
func DoAndReturn(f func() int, fallback int) (result int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Unexpected panic was caught:", err)
			result = fallback
		}
	}()
	return f()
}
