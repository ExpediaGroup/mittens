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
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	host       string
}

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
	return Client{httpClient: client, host: host}
}

// Send http request and returns error also if response code is NOT 2XX
func (c Client) Request(method, path string, headers map[string]string, requestBody *string) error {

	var body io.Reader
	if requestBody != nil {
		body = bytes.NewBuffer([]byte(*requestBody))
	}

	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.host, "/"), strings.TrimLeft(path, "/"))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("new http request: %s %s: %v", method, url, err)
	}

	for k, v := range headers {
		if strings.EqualFold(k, "Host") {
			req.Host = v
		}
		req.Header.Add(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}

	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body) // try to read response body as well to give user more info why request failed
		return fmt.Errorf("%s %s returned %d %s, expected 2xx",
			method, url, resp.StatusCode, strings.TrimSuffix(string(b), "\n"))
	}

	if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
		return fmt.Errorf("read response body: %s %s: %v", method, url, err)
	}
	return nil
}
