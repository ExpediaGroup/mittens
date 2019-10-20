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
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ToHeaders(t *testing.T) {

	headersFlag := []string{
		"Content-Type: application/json",
		"Host: localhost",
	}

	headers := toHeaders(headersFlag)

	assert.Equal(t, 2, len(headers))
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, "localhost", headers["Host"])
}

func Test_HeadersWithoutSeparatorToHeaders(t *testing.T) {

	headersFlag := []string{
		"No separator",
	}

	headers := toHeaders(headersFlag)

	assert.Equal(t, 1, len(headers))
	value, ok := headers["No separator"]
	assert.True(t, ok)
	assert.Equal(t, "", value)
}

func Test_HeadersWithMultipleSeparatorsToHeaders(t *testing.T) {

	headersFlag := []string{
		"Cookie: some:strange:cookie",
	}

	headers := toHeaders(headersFlag)

	assert.Equal(t, 1, len(headers))
	assert.Equal(t, "some:strange:cookie", headers["Cookie"])
}
