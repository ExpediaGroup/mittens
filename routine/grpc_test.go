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

package routine

import (
	"fmt"
	"testing"
)

type grpcRequestTestCase struct {
	input          string
	expectedOutput grpcRequest
	expectedErr    error
}

func TestCreateGrpcRequest(t *testing.T) {
	testCases := []grpcRequestTestCase{
		{
			":{}",
			grpcRequest{},
			fmt.Errorf("unknown grpc request format: %s", ":{}"),
		},
		{
			"service/method:",
			grpcRequest{},
			fmt.Errorf("unknown grpc request format: %s", "service/method:"),
		},
		{
			"service/method:{}",
			grpcRequest{
				methodName: "service/method",
				message:    "{}",
			},
			nil,
		},
		{
			"service/method:{\"id\": \"1234\", \"text\": \"test\"}",
			grpcRequest{
				methodName: "service/method",
				message:    "{\"id\": \"1234\", \"text\": \"test\"}",
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		actual, err := createGrpcRequest(testCase.input)

		if (testCase.expectedErr != nil || err != nil) && (testCase.expectedErr.Error() != err.Error()) {
			t.Fatalf("Expected error\n%v\nbut got\n%v", testCase.expectedErr, err)
		} else if actual != testCase.expectedOutput {
			t.Fatalf("Expected output\n%s\nbut got\n%s", testCase.expectedOutput, actual)
		}
	}
}
