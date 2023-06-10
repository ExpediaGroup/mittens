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

package grpc

import (
	"bytes"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"mittens/internal/pkg/placeholders"
	"mittens/internal/pkg/response"
	"os"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Client represents a gRPC client.
type Client struct {
	host                string
	insecure            bool
	timeoutMilliseconds int
	connClose           func() error
	conn                *grpc.ClientConn
	descriptorSource    grpcurl.DescriptorSource
}

// eventHandler is a custom event handler with the option to enable/disable logging of responses.
type eventHandler struct {
	grpcurl.InvocationEventHandler
	logResponses bool
}

// NewClient returns a gRPC client.
func NewClient(host string, insecure bool, timeoutMilliseconds int) Client {
	return Client{host: host, insecure: insecure, connClose: func() error { return nil }, timeoutMilliseconds: timeoutMilliseconds}
}

// Connect attempts to establish a connection with a gRPC server.
func (c *Client) Connect(headers []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeoutMilliseconds)*time.Millisecond)

	headersMetadata := grpcurl.MetadataFromHeaders(headers)
	contextWithMetadata := metadata.NewOutgoingContext(ctx, headersMetadata)

	// grpc.WithReturnConnectionError() is EXPERIMENTAL and may be changed or removed in a later release
	// Added to provide more information if a connection error occurs
	dialOptions := []grpc.DialOption{grpc.WithBlock(), grpc.WithReturnConnectionError()}
	if c.insecure {
		log.Print("ignoring gRPC server SSL/TLS authentication")
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		log.Print("using gRPC server SSL/TLS authentication")
		tlsConf, err := grpcurl.ClientTLSConfig(c.insecure, "", "", "")
		if err != nil {
			return fmt.Errorf("failed to create TLS config: %v", err)
		}
		var creds = credentials.NewTLS(tlsConf)
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	}

	grpc.WithReturnConnectionError()
	conn, err := grpc.DialContext(ctx, c.host, dialOptions...)
	if err != nil {
		return fmt.Errorf("gRPC dial: %v", err)
	}

	reflectionClient := grpcreflect.NewClientV1Alpha(contextWithMetadata, reflectpb.NewServerReflectionClient(conn))
	descriptorSource := grpcurl.DescriptorSourceFromServer(contextWithMetadata, reflectionClient)

	log.Print("gRPC client connected")
	c.conn = conn
	c.connClose = func() error { cancel(); return conn.Close() }
	c.descriptorSource = descriptorSource
	return nil
}

// SendRequest sends a request to the gRPC server and wraps useful information into a Response object.
// Note that the message cannot be null. Even if there is no message to be sent this needs to be set to an empty string.
func (c *Client) SendRequest(serviceMethod string, message string, headers []string, logResponses bool) response.Response {
	const respType = "grpc"
	in := bytes.NewBufferString(message)

	// TODO - create generic parser and formatter for any request, can we use text parser/formatter?
	requestParser, formatter, err := grpcurl.RequestParserAndFormatter("json", c.descriptorSource, in, grpcurl.FormatOptions{})
	if err != nil {
		log.Printf("Cannot construct request parser and formatter for json")
		// FIXME FATAL
		return response.Response{Duration: time.Duration(0), Err: err, Type: respType}
	}

	delegate := &grpcurl.DefaultEventHandler{
		Out:       os.Stdout,
		Formatter: formatter,
	}

	loggingEventHandler := eventHandler{InvocationEventHandler: delegate, logResponses: logResponses}
	startTime := time.Now()

	// Interpolate
	interpolatedHeaders := make([]string, len(headers))
	for i, header := range headers {
		interpolatedHeaders[i] = placeholders.InterpolatePlaceholders(header)
	}

	err = grpcurl.InvokeRPC(context.Background(), c.descriptorSource, c.conn, serviceMethod, interpolatedHeaders, loggingEventHandler, requestParser.Next)
	endTime := time.Now()
	if err != nil {
		log.Printf("grpc response error: %s", err)
		return response.Response{Duration: endTime.Sub(startTime), Err: nil, Type: respType}
	}
	return response.Response{Duration: endTime.Sub(startTime), Err: nil, Type: respType}
}

// OnReceiveResponse overrides the default method and allows enabling/disabling logging of responses.
func (h eventHandler) OnReceiveResponse(msg proto.Message) {
	if h.logResponses {
		h.InvocationEventHandler.OnReceiveResponse(msg)
	}
}

// Close calling close on a client that has not established connection does not return an error.
func (c Client) Close() error {
	log.Print("Closing gRPC client connection")
	return c.connClose()
}
