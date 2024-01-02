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
	"mittens/internal/pkg/placeholders"
	"mittens/internal/pkg/response"
	"mittens/internal/pkg/util"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

// Client is a wrapper for the HTTP Client which includes a host.
type Client struct {
	httpClient *http.Client
	host       string
}

type ProtocolType string

const (
	HTTP1 ProtocolType = "h1"
	HTTP2 ProtocolType = "h2"
	H2C   ProtocolType = "h2c"
)

// NewClient creates a new HTTP client for a given host.
// If insecure is true, the client will not verify the server's certificate chain and host name.
func NewClient(host string, insecure bool, timeoutMilliseconds int, protocol ProtocolType) Client {
	client := &http.Client{
		Timeout: time.Duration(timeoutMilliseconds) * time.Millisecond,
	}

	switch protocol {
	case HTTP2:
		client.Transport = &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
	case H2C:
		client.Transport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
	default:
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
	}

	return Client{httpClient: client, host: strings.TrimRight(host, "/")}
}

// SendRequest sends a request to the HTTP server and wraps useful information into a Response object.
func (c Client) SendRequest(method, path string, headers []string, requestBody *string) response.Response {
	const respType = "http"
	var body io.Reader
	if requestBody != nil {
		body = bytes.NewBufferString(*requestBody)
	}

	url := fmt.Sprintf("%s/%s", c.host, strings.TrimLeft(path, "/"))
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		log.Printf("Failed to create request: %s %s: %v", method, url, err)
		return response.Response{Duration: time.Duration(0), Err: err, Type: respType}
	}

	headersMap := util.ToHeaders(headers)
	for k, v := range headersMap {
		if strings.EqualFold(k, "Host") {
			req.Host = v
		}

		// interpolate headers (just the values, not the keys)
		interpolatedHeaderValue := placeholders.InterpolatePlaceholders(v)

		req.Header.Add(k, interpolatedHeaderValue)
	}

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	endTime := time.Now()
	if err != nil {
		return response.Response{Duration: endTime.Sub(startTime), Err: err, Type: respType}
	}
	defer resp.Body.Close()

	if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
		return response.Response{Duration: endTime.Sub(startTime), Err: err, Type: respType, StatusCode: resp.StatusCode}
	}
	return response.Response{Duration: endTime.Sub(startTime), Err: nil, Type: respType, StatusCode: resp.StatusCode}
}
