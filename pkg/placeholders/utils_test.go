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

package placeholders

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttp_DateInterpolation(t *testing.T) {
	input := `post:/db_{$currentDate}:{"date": "{$currentDate|days+5,months+2,years-1,format=yyyy-MM-dd}"}`
	output := InterpolatePlaceholders(input)

	dateToday := time.Now().Format("2006-01-02")                        // today
	dateWithOffset := time.Now().AddDate(-1, 2, 5).Format("2006-01-02") // today -1 year, +2 months, +5 days
	assert.Equal(t, fmt.Sprintf(`post:/db_%s:{"date": "%s"}`, dateToday, dateWithOffset), output)
}

func TestHttp_TimestampInterpolation(t *testing.T) {
	input := `post:/path_{$currentTimestamp}:{"body": "{$currentTimestamp}"}`
	output := InterpolatePlaceholders(input)

	assert.Equal(t, len(output), 50) //  13 numbers for timestamp + 13 for timestamp + 24 for rest of text
}

func TestHttp_MultipleInterpolation(t *testing.T) {
	input := `post:/path_{$range|min=1,max=2}_{$random|foo,bar}:{"body": "{$random|foo,bar} {$range|min=1,max=2}"}`
	output := InterpolatePlaceholders(input)

	var outputRegex = regexp.MustCompile("/path_\\d_(foo|bar):{\"body\": \"(foo|bar) \\d\"}")
	matchOutput := outputRegex.MatchString(output)

	assert.True(t, matchOutput)
}

func TestHttp_RangeInterpolation(t *testing.T) {
	input := `post:/path_{$range|min=1,max=2}:{"body": "{$range|min=1,max=2}"}`
	output := InterpolatePlaceholders(input)

	var outputRegex = regexp.MustCompile("/path_\\d:{\"body\": \"\\d\"}")
	matchOutput := outputRegex.MatchString(output)

	assert.True(t, matchOutput)
}

func TestHttp_InvalidRangeInterpolation(t *testing.T) {
	input := `post:/path_{$range|min=2,max=1}:{"body": "{$range|min=2,max=1}"}`
	output := InterpolatePlaceholders(input)

	// will not action on invalid ranges
	assert.Equal(t, output, "post:/path_{$range|min=2,max=1}:{\"body\": \"{$range|min=2,max=1}\"}")
}

func TestHttp_RandomElementInterpolation(t *testing.T) {
	input := `{$random|fo-o,b_ar}:{"body": "{$random|fo-o,b_ar}"}`
	output := InterpolatePlaceholders(input)

	var elementsRegex = regexp.MustCompile("(fo-o|b_ar)")
	matchOutput := elementsRegex.MatchString(output)

	assert.True(t, matchOutput)
}
