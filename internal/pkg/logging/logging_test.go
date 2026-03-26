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

package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup_TextFormat(t *testing.T) {
	err := Setup("info", "text", "iso8601")
	assert.NoError(t, err)
}

func TestSetup_JsonFormat(t *testing.T) {
	err := Setup("info", "json", "iso8601")
	assert.NoError(t, err)
}

func TestSetup_AllLogLevels(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "warning", "error"} {
		err := Setup(level, "text", "iso8601")
		assert.NoError(t, err, "level %s should be valid", level)
	}
}

func TestSetup_InvalidLogLevel(t *testing.T) {
	err := Setup("invalid", "text", "iso8601")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported log level")
}

func TestSetup_InvalidLogFormat(t *testing.T) {
	err := Setup("info", "xml", "iso8601")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported log format")
}

func TestSetup_AllDateTimeFormats(t *testing.T) {
	for _, format := range []string{"iso8601", "rfc3339", "rfc3339nano", "epoch"} {
		err := Setup("info", "text", format)
		assert.NoError(t, err, "datetime format %s should be valid", format)
	}
}
