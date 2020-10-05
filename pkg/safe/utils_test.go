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
	"github.com/stretchr/testify/assert"
	"testing"
)

const RESULT = 1
const FALLBACK = 2

func TestPanicIsCaught(t *testing.T) {
	assert.NotPanics(t, func() {
		Do(func() {
			panic("test panic")
		})
	})

	assert.NotPanics(t, func() {
		DoAndReturn(func() int {
			panic("test panic")
			return RESULT
		}, FALLBACK)
	})
}

func TestOriginalResultReturnedWhenNoPanic(t *testing.T) {
	actual := DoAndReturn(func() int {
		return RESULT
	}, FALLBACK)

	assert.Equal(t, RESULT, actual)
}

func TestFallbackResultReturnedWhenPanic(t *testing.T) {
	actual := DoAndReturn(func() int {
		panic("test panic")
		return RESULT
	}, FALLBACK)

	assert.Equal(t, FALLBACK, actual)
}

func TestHasPanicked(t *testing.T) {
	Do(func() { panic("test panic") })

	assert.True(t, HasPanicked())
}
