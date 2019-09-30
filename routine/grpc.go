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

// Given a string in the form method:content it creates a gRPC request
// that will be used to call a service/method using grpcurl.

package routine

import (
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var grpcContentRegex = regexp.MustCompile("^(.+?):(.+?)$")

type grpcRequest struct {
	methodName string
	message    string
}

func createGrpcRequest(content string) (grpcRequest, error) {
	var request grpcRequest
	var err error = nil

	if grpcContentRegex.MatchString(content) {
		matched := grpcContentRegex.FindStringSubmatch(content)
		request = grpcRequest{
			methodName: matched[1],
			message:    matched[2],
		}
	} else {
		err = fmt.Errorf("unknown grpc request format: %s", content)
	}

	return request, err
}

func sendGrpcRequest(descriptorSource grpcurl.DescriptorSource, clientConnection *grpc.ClientConn, grpcHeaders []string, warmupReq grpcRequest) error {
	in := strings.NewReader(warmupReq.message)
	jsonRequestParser, jsonFormatter, err := grpcurl.RequestParserAndFormatterFor(grpcurl.Format("json"), descriptorSource, false, false, in)
	if err != nil {
		log.Printf("Failed to construct jsonFormatter for json with error %v\n", err)
		return err
	}
	loggingEventHandler := grpcurl.NewDefaultEventHandler(os.Stdout, descriptorSource, jsonFormatter, false)

	log.Printf("calling: RPC %s with message: %s\n", warmupReq.methodName, warmupReq.message)
	err = grpcurl.InvokeRPC(context.Background(), descriptorSource, clientConnection, warmupReq.methodName, grpcHeaders, loggingEventHandler, jsonRequestParser.Next)
	if err != nil {
		log.Printf("gRPC call failed with error %v\n", err)
		time.Sleep(time.Second)
	}
	return err
}

// Creates a gRPC client that uses reflection to call the gRCP server
func createGrpcClient(grpcHeaders StringArray, targetHost string, targetGrpcPort int) (*grpc.ClientConn, grpcurl.DescriptorSource, error) {
	emptyContext := context.Background()
	headersMetadata := grpcurl.MetadataFromHeaders(grpcHeaders)
	contextWithMetadata := metadata.NewOutgoingContext(emptyContext, headersMetadata)

	clientConnection, err := grpc.DialContext(emptyContext, fmt.Sprintf("%s:%d", targetHost, targetGrpcPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("gRPC client creation failed with error %v\n", err)
		return nil, nil, err
	}

	reflectionClient := grpcreflect.NewClient(contextWithMetadata, reflectpb.NewServerReflectionClient(clientConnection))
	descriptorSource := grpcurl.DescriptorSourceFromServer(contextWithMetadata, reflectionClient)
	return clientConnection, descriptorSource, nil
}
