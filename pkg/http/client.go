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
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mittens/pkg/response"
	"net/http"
	"strings"
	"time"
)

// Client is a wrapper for the httpClient which includes a host
type Client struct {
	httpClient *http.Client
	host       string
}

// NewClient creates a new client for a given host. if insecure is true,
// client won't verify the server's certificate chain and host name
func NewClient(host string, insecure bool) Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if insecure {
		log.Printf("http client: insecure skip verify is set to true")
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return Client{httpClient: client, host: strings.TrimRight(host, "/")}
}

// Request sends an http request and wraps useful info into a Response object
func (c Client) Request(method, path string, headers map[string]string, requestBody *string) response.Response {
	var body io.Reader
	if requestBody != nil {
		body = bytes.NewBufferString(*requestBody)
	}

	url := fmt.Sprintf("%s/%s", c.host, strings.TrimLeft(path, "/"))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("new http request: %s %s: %v", method, url, err)
		return response.Response{Duration: time.Duration(0), Err: err, ClientError: false, RequestSent: false, Type: "http"}
	}

	for k, v := range headers {
		if strings.EqualFold(k, "Host") {
			req.Host = v
		}
		req.Header.Add(k, v)
	}

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	endTime := time.Now()
	if err != nil {
		log.Printf("http request: %v", err)
		return response.Response{Duration: endTime.Sub(startTime), Err: err, ClientError: false, RequestSent: false, Type: "http"}
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return response.Response{Duration: endTime.Sub(startTime), Err: fmt.Errorf("statusCode: %d\n", resp.StatusCode), ClientError: resp.StatusCode/100 == 4, RequestSent: true,
			Type: "http"}
	}

	if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
		log.Printf("read response body: %s %s: %v", method, url, err)
		return response.Response{Duration: endTime.Sub(startTime), Err: err, ClientError: false, RequestSent: true, Type: "http"}
	}
	return response.Response{Duration: endTime.Sub(startTime), Err: nil, ClientError: false, RequestSent: true, Type: "http"}
}
