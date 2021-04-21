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
	"log"
	"mittens/pkg/response"
	"os"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Client represents a gRPC client.
type Client struct {
	host             string
	timeoutSeconds   int
	insecure         bool
	connClose        func() error
	conn             *grpc.ClientConn
	descriptorSource grpcurl.DescriptorSource
}

// then define eventHandler like so
type eventHandler struct {
	grpcurl.InvocationEventHandler
	logResponses bool
}

// NewClient returns a gRPC client.
func NewClient(host string, insecure bool, timeoutSeconds int) Client {
	return Client{host: host, timeoutSeconds: timeoutSeconds, insecure: insecure, connClose: func() error { return nil }}
}

// Connect attempts to establish a connection with a gRPC server.
func (c *Client) Connect(headers []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeoutSeconds)*time.Second)

	headersMetadata := grpcurl.MetadataFromHeaders(headers)
	contextWithMetadata := metadata.NewOutgoingContext(ctx, headersMetadata)

	dialOptions := []grpc.DialOption{grpc.WithBlock()}
	if c.insecure {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	conn, err := grpc.DialContext(ctx, c.host, dialOptions...)
	if err != nil {
		return fmt.Errorf("gRPC dial: %v", err)
	}

	reflectionClient := grpcreflect.NewClient(contextWithMetadata, reflectpb.NewServerReflectionClient(conn))
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

	//  use this to create your handler
	delegate := &grpcurl.DefaultEventHandler{
		Out:       os.Stdout,
		Formatter: formatter,
	}

	loggingEventHandler := eventHandler{InvocationEventHandler: delegate, logResponses: logResponses}
	startTime := time.Now()
	err = grpcurl.InvokeRPC(context.Background(), c.descriptorSource, c.conn, serviceMethod, headers, loggingEventHandler, requestParser.Next)
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
