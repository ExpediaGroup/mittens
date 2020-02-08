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
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const lettersAndNumbers = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers = "0123456789"

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

var templateRegex = regexp.MustCompile("{(chars([-]\\d+)?|numbers([-]\\d+)?|today([+-]\\d+)?|tomorrow)}")
var templatePlusMinusRegex = regexp.MustCompile("[+-]\\d+")
var templateOffsetRegex = regexp.MustCompile("\\d+")

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

func interpolatePlaceholders(source string) string {
	return templateRegex.ReplaceAllStringFunc(source, func(templateString string) string {

		if strings.Contains(templateString, "today") || strings.Contains(templateString, "tomorrow") {
			return dateElements(templateString)
		} else {
			return randomElements(templateString)
		}
	})
}

func randomElements(source string) string {

	if strings.Contains(source, "numbers") {
		length, _ := strconv.Atoi(templateOffsetRegex.FindString(source))
		return randomNumbers(length)
	} else {
		length, _ := strconv.Atoi(templateOffsetRegex.FindString(source))
		return randomLetters(length)
	}
}

func randomNumbers(n int) string {
	output := make([]byte, n)
	for i := range output {
		output[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(output)
}

func randomLetters(n int) string {
	output := make([]byte, n)
	for i := range output {
		output[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}
	return string(output)
}

func dateElements(source string) string {
	offsetDays := 0

	if source == "{today}" {
		offsetDays = 0
	} else if source == "{tomorrow}" {
		offsetDays = 1
	} else if extractedOffset := templatePlusMinusRegex.FindString(source); len(extractedOffset) > 0 {
		offsetDays, _ = strconv.Atoi(extractedOffset)
	}

	// the date below is how the golang date formatter works. it's used for the formatting. it's not what is actually going to be displayed
	return time.Now().Add(time.Duration(offsetDays) * 24 * time.Hour).Format("2006-01-02")
}
