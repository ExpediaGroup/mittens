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

package http

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Method string
	Path   string
	Body   *string
}

var allowedHttpMethods = map[string]interface{}{
	"GET":     nil,
	"HEAD":    nil,
	"POST":    nil,
	"PUT":     nil,
	"PATCH":   nil,
	"DELETE":  nil,
	"CONNECT": nil,
	"OPTIONS": nil,
	"TRACE":   nil,
}

var templatePlaceholderRegex = regexp.MustCompile("{{(\\w+([\\|+][\\w+-=,]+)?|(\\w+))}}")
var templateRangeRegex = regexp.MustCompile("{{range\\|min=(?P<Min>\\d+),max=(?P<Max>\\d+)}}")
var templateElementsRegex = regexp.MustCompile("{{random\\|(?P<Elements>[,\\w-]+)}}")
var templateDatesRegex = regexp.MustCompile("{{currentDate(?:\\|(?:days(?P<Days>[+-]\\d+))*(?:[,]*months(?P<Months>[+-]\\d+))*(?:[,]*years(?P<Years>[+-]\\d+))*)*}}") // we did it \o/ 100+ chars regex

func ToHttpRequest(requestFlag string) (Request, error) {
	parts := strings.SplitN(requestFlag, ":", 3)
	if len(parts) < 2 {
		return Request{}, fmt.Errorf("invalid request flag: %s, expected format <http-method>:<path>[:body]", requestFlag)
	}

	method := strings.ToUpper(parts[0])
	_, ok := allowedHttpMethods[method]
	if !ok {
		return Request{}, fmt.Errorf("invalid request flag: %s, method %s is not supported", requestFlag, method)
	}

	// <method>:<path>
	if len(parts) == 2 {
		path := interpolatePlaceholders(parts[1])

		return Request{
			Method: method,
			Path:   path,
			Body:   nil,
		}, nil
	}

	path := interpolatePlaceholders(parts[1])
	var body = interpolatePlaceholders(parts[2])

	return Request{
		Method: method,
		Path:   path,
		Body:   &body,
	}, nil
}

func dateElements(source string) string {
	r := templateDatesRegex.FindStringSubmatch(source)

	if r == nil {
		return source
	}
	days := r[1]
	months := r[2]
	years := r[3]

	offsetDays, _ := strconv.Atoi(days)
	offsetMonths, _ := strconv.Atoi(months)
	offsetYears, _ := strconv.Atoi(years)

	// the date below is how the golang date formatter works. it's used for the formatting. it's not what is actually going to be displayed
	return time.Now().AddDate(offsetYears, offsetMonths, offsetDays).Format("2006-01-02")
}

func timestampElements() string {
	epoch := time.Now().UnixNano() / 1000000

	return strconv.FormatInt(epoch, 10)
}

func randomElements(source string) string {
	r := templateElementsRegex.FindStringSubmatch(source)

	if r == nil {
		return source
	}

	s := strings.Split(r[1], ",")
	number := rand.Intn(len(s))

	return s[number]
}

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

func interpolatePlaceholders(source string) string {
	return templatePlaceholderRegex.ReplaceAllStringFunc(source, func(templateString string) string {
		rand.Seed(time.Now().UnixNano())

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
