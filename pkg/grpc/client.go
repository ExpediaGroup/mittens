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
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Client struct {
	host             string
	timeoutSeconds   int
	insecure         bool
	grpcConnectOnce  *sync.Once
	connClose        func() error
	conn             *grpc.ClientConn
	descriptorSource grpcurl.DescriptorSource
}

func NewClient(host string, insecure bool, timeoutSeconds int) Client {
	return Client{host: host, timeoutSeconds: timeoutSeconds, grpcConnectOnce: new(sync.Once), insecure: insecure, connClose: func() error { return nil }}
}

func (c *Client) connect(headers []string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeoutSeconds)*time.Second)
	headersMetadata := grpcurl.MetadataFromHeaders(headers)
	contextWithMetadata := metadata.NewOutgoingContext(ctx, headersMetadata)

	dialOptions := []grpc.DialOption{grpc.WithBlock()}
	if c.insecure {
		log.Print("grpc client insecure")
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	log.Printf("grpc client connecting to %s", c.host)
	conn, err := grpc.Dial(c.host, dialOptions...)
	if err != nil {
		return fmt.Errorf("grpc dial: %v", err)
	}

	reflectionClient := grpcreflect.NewClient(contextWithMetadata, reflectpb.NewServerReflectionClient(conn))
	descriptorSource := grpcurl.DescriptorSourceFromServer(contextWithMetadata, reflectionClient)

	log.Print("grpc client connected")
	c.conn = conn
	c.connClose = func() error { cancel(); return conn.Close() }
	c.descriptorSource = descriptorSource
	return nil
}

func (c *Client) Request(serviceMethod string, headers []string, body []byte) error {

	var connErr error
	c.grpcConnectOnce.Do(func() {
		connErr = c.connect(headers)
	})

	if connErr != nil {
		return fmt.Errorf("grpc client connect: %v", connErr)
	}

	var in io.Reader
	if len(body) != 0 {
		in = bytes.NewReader(body)
	}

	// TODO - create generic parser and formatter for any request, can we use text parser/formatter?
	requestParser, formatter, err := grpcurl.RequestParserAndFormatterFor(grpcurl.Format("json"), c.descriptorSource, false, false, in)
	if err != nil {
		return fmt.Errorf("cannot construct request parser and formatter for json")
	}
	loggingEventHandler := grpcurl.NewDefaultEventHandler(os.Stdout, c.descriptorSource, formatter, false)
	return grpcurl.InvokeRPC(context.Background(), c.descriptorSource, c.conn, serviceMethod, headers, loggingEventHandler, requestParser.Next)
}

// Calling close on client that has not established connection does not return error
func (c Client) Close() error {
	log.Print("closing grpc client connection")
	return c.connClose()
}