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
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"mittens/internal/pkg/placeholders"
	"strings"

	"github.com/andybalholm/brotli"
)

// Request represents an HTTP request.
type Request struct {
	Method  string
	Headers map[string]string
	Path    string
	Body    *string
}

type CompressionType string

const (
	COMPRESSION_NONE    CompressionType = ""
	COMPRESSION_GZIP    CompressionType = "gzip"
	COMPRESSION_BROTLI  CompressionType = "brotli"
	COMPRESSION_DEFLATE CompressionType = "deflate"
)

var allowedHTTPMethods = map[string]interface{}{
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

// ToHTTPRequest parses an HTTP request which is in a string format and stores it in a struct.
func ToHTTPRequest(requestString string, compression CompressionType) (Request, error) {
	parts := strings.SplitN(requestString, ":", 3)
	if len(parts) < 2 {
		return Request{}, fmt.Errorf("invalid request flag: %s, expected format <http-method>:<path>[:body]", requestString)
	}

	method := strings.ToUpper(parts[0])
	_, ok := allowedHTTPMethods[method]
	if !ok {
		return Request{}, fmt.Errorf("invalid request flag: %s, method %s is not supported", requestString, method)
	}

	// <method>:<path>
	if len(parts) == 2 {
		path := placeholders.InterpolatePlaceholders(parts[1])

		return Request{
			Method: method,
			Path:   path,
			Body:   nil,
		}, nil
	}

	path := placeholders.InterpolatePlaceholders(parts[1])
	// the body of the request can either be inlined, or come from a file
	rawBody, err := placeholders.GetBodyFromFileOrInlined(parts[2])
	if err != nil {
		return Request{}, fmt.Errorf("unable to parse body for request: %s", parts[2])
	}
	var body = placeholders.InterpolatePlaceholders(*rawBody)

	var reader io.Reader
	switch compression {
	case COMPRESSION_GZIP:
		reader, err = compressGzip([]byte(body))
	case COMPRESSION_BROTLI:
		reader, err = compressBrotli([]byte(body))
	case COMPRESSION_DEFLATE:
		reader, err = compressFlate([]byte(body))
	default:
		reader = bytes.NewBufferString(body)
	}

	if err != nil {
		return Request{}, fmt.Errorf("unable to compress body for request: %s", parts[2])
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	compressedBody := buf.String()

	headers := make(map[string]string)
	if compression != COMPRESSION_NONE {
		encoding := ""
		switch compression {
		case COMPRESSION_GZIP:
			encoding = "gzip"
		case COMPRESSION_BROTLI:
			encoding = "br"
		case COMPRESSION_DEFLATE:
			encoding = "deflate"
		}
		headers["Content-Encoding"] = encoding
	}

	return Request{
		Method:  method,
		Headers: headers,
		Path:    path,
		Body:    &compressedBody,
	}, nil
}

func compressGzip(data []byte) (io.Reader, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data); if err != nil {
		return nil, err
	} 
    err = w.Close(); if err != nil {
		return nil, err
	} 
	return &b, nil
}

func compressFlate(data []byte) (io.Reader, error) {
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, 9)
	_, err := w.Write(data); if err != nil {
		return nil, err
	} 
    err = w.Close(); if err != nil {
		return nil, err
	} 
	return &b, nil
}

func compressBrotli(data []byte) (io.Reader, error) {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	_, err := w.Write(data); if err != nil {
		return nil, err
	} 
    err = w.Close(); if err != nil {
		return nil, err
	} 
	return &b, nil
}
