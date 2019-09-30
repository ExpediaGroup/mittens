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

// Given a string in the form method:url:body it creates an HTTP (REST) request
// that will be used to call a REST endpoint.

package routine

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// request is in the form method:url:body where method is one of [get, post, or put],
// url is the url where the request will be sent, and body is a (optional) properly escaped JSON-formatted string.
var httpContentRegex = regexp.MustCompile("^(get|post|put):(/.*?)(:(.*?))?$")
var todayTemplateRegex = regexp.MustCompile("{(today([+-]\\d+)?|tomorrow)}")
var todayTemplatePlusMinusRegex = regexp.MustCompile("[+-]\\d+")

type httpRequest struct {
	method string
	url    string
	body   string
}

// Parses the given string and extracts the method, url, and body of the HTTP request.
func createHttpRequest(content string) (httpRequest, error) {
	var request httpRequest
	var err error = nil

	if httpContentRegex.MatchString(content) {
		matched := httpContentRegex.FindStringSubmatch(content)

		body := ""
		if len(matched) == 5 {
			body = matched[4]
		}

		request = httpRequest{
			method: strings.ToUpper(matched[1]),
			url:    interpolateDates(matched[2]),
			body:   interpolateDates(body),
		}
	} else {
		err = fmt.Errorf("unknown http request format: %s", content)
	}

	return request, err
}

// Creates a HTTP client with a given timeout and SSL checks disabled.
func createHttpClient(timeoutSeconds int) *http.Client {
	// disable ssl checks
	transportNoCerts := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// timeDuration uses nano seconds...
	timeout := time.Duration(timeoutSeconds) * 1000 * 1000 * 1000
	// create http client
	client := &http.Client{
		Transport: transportNoCerts,
		Timeout:   timeout,
	}
	return client
}

// Provides date templating; given a string which contains {today} or {today+/-n}
// it returns the date in YYYY-MM-DD format.
func interpolateDates(source string) string {
	return todayTemplateRegex.ReplaceAllStringFunc(source, func(templateString string) string {
		offsetDays := 0

		if templateString == "{tomorrow}" {
			offsetDays = 1
		} else if extractedOffset := todayTemplatePlusMinusRegex.FindString(templateString); len(extractedOffset) > 0 {
			offsetDays, _ = strconv.Atoi(extractedOffset)
		}

		// the date below is how the golang date formatter works. it's used for the formatting. it's not what is actually going to be displayed
		return time.Now().Add(time.Duration(offsetDays) * 24 * time.Hour).Format("2006-01-02")
	})
}

// Sends an HTTP request to the given endpoint using the provided HTTP client.
func sendHttpRequest(client *http.Client, method string, url string, httpHeaders *map[string]string, body string) (int, error) {
	escapedUrl := escapeQueryStringValues(url)

	if len(body) > 0 {
		log.Printf("calling: %s %s with body: %s\n", method, escapedUrl, body)
	} else {
		log.Printf("calling: %s %v\n", method, escapedUrl)
	}

	req, _ := http.NewRequest(method, escapedUrl, strings.NewReader(body))

	for k, v := range *httpHeaders {
		if strings.EqualFold("Host", k) {
			req.Host = v
		} else {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)

	if err == nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %s\n", err)
			}
		}()

		message, _ := ioutil.ReadAll(resp.Body)
		if len(message) > 0 {
			log.Printf("returned: %v with body: %s\n", resp.StatusCode, message)
		} else {
			log.Printf("returned: %v\n", resp.StatusCode)
		}

		return resp.StatusCode, err
	} else {
		log.Printf("Error: %s\n", err)
		return 0, err
	}
}

// Receives an URL as a string, escapes the query params and returns it back as a string.
func escapeQueryStringValues(fullUrlStr string) string {
	fullUrl, err := url.Parse(fullUrlStr)

	if err != nil {
		log.Printf("Error parsing URL %s: %s ", fullUrlStr, err)
		return fullUrlStr
	}

	values, err := url.ParseQuery(fullUrl.RawQuery)

	if err != nil {
		log.Printf("Error parsing query string %s: %s ", fullUrl.RawQuery, err)
		return fullUrlStr
	}

	fullUrl.RawQuery = values.Encode()
	return fullUrl.String()
}
