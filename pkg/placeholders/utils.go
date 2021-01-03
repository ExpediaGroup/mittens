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
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// anything that starts with {$, followed by any word character, and optionally followed by a modifier identifier | and the modifiers that can contain word chars + - = and ,
var templatePlaceholderRegex = regexp.MustCompile("{\\$(\\w+(?:[\\|(?:[\\w+-=,]+)]*)}")
var templateRangeRegex = regexp.MustCompile("{\\$range\\|min=(?P<Min>\\d+),max=(?P<Max>\\d+)}")
var templateElementsRegex = regexp.MustCompile("{\\$random\\|(?P<Elements>[,\\w-]+)}")
var templateDatesRegex = regexp.MustCompile("{\\$currentDate(?:\\|(?:days(?P<Days>[+-]\\d+))*(?:[,]*months(?P<Months>[+-]\\d+))*(?:[,]*years(?P<Years>[+-]\\d+))*(?:[,]*format=(?P<Format>[yMd|,/-]+))*)*}")

// dateElements replaces date placeholders with the actual dates. It supports offsets for days, months, and years.
func dateElements(source string) string {
	r := templateDatesRegex.FindStringSubmatch(source)

	if r == nil {
		return source
	}
	days := r[1]
	months := r[2]
	years := r[3]
	format := r[4]

	offsetDays, _ := strconv.Atoi(days)
	offsetMonths, _ := strconv.Atoi(months)
	offsetYears, _ := strconv.Atoi(years)

	if format == "" {
		// If no format override is specified, we default to ISO 8601, or YYYY MM DD
		format = "2006-01-02"
	} else {
		format = strings.ReplaceAll(format, "yyyy", "2006")
		format = strings.ReplaceAll(format, "yy", "06")
		format = strings.ReplaceAll(format, "MMM", "Jan")
		format = strings.ReplaceAll(format, "MM", "01")
		format = strings.ReplaceAll(format, "dd", "02")
		format = strings.ReplaceAll(format, "d", "2")
	}

	// the date below is how the golang date formatter works. it's used for the formatting. it's not what is actually going to be displayed
	return time.Now().AddDate(offsetYears, offsetMonths, offsetDays).Format(format)
}

// timestampElements returns the current time from Unix epoch in milliseconds.
func timestampElements() string {
	epoch := time.Now().UnixNano() / 1000000

	return strconv.FormatInt(epoch, 10)
}

// randomElements replaces random element placeholders with elements which are randomly selected from the provided list.
func randomElements(source string) string {
	r := templateElementsRegex.FindStringSubmatch(source)

	if r == nil {
		return source
	}

	s := strings.Split(r[1], ",")
	number := rand.Intn(len(s))

	return s[number]
}

// rangeElements replaces range element placeholders with random integers within the specified range.
func rangeElements(source string) string {
	r := templateRangeRegex.FindStringSubmatch(source)
	if r == nil {
		return source
	}

	min, _ := strconv.Atoi(r[1])
	max, _ := strconv.Atoi(r[2])

	if min > max {
		log.Printf("Invalid range. min > max")
		return source
	}

	number := rand.Intn(max-min+1) + min

	return strconv.Itoa(number)
}

// InterpolatePlaceholders scans a string and replaces placeholders with actual values.
// At the moment this supports; dates, timestamps, random values from a list, and random integers.
func InterpolatePlaceholders(source string) string {
	return templatePlaceholderRegex.ReplaceAllStringFunc(source, func(templateString string) string {

		if strings.Contains(templateString, "currentDate") {
			return dateElements(templateString)
		} else if strings.Contains(templateString, "currentTimestamp") {
			return timestampElements()
		} else if strings.Contains(templateString, "random") {
			return randomElements(templateString)
		} else if strings.Contains(templateString, "range") {
			return rangeElements(templateString)
		} else {
			return source
		}
	})
}
