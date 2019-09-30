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
	"time"
)

type httpRequestTestCase struct {
	input          string
	expectedOutput httpRequest
	expectedErr    error
}

func TestCreateHttpRequest(t *testing.T) {
	testCases := []httpRequestTestCase{
		{
			"foo:/bar:",
			httpRequest{},
			fmt.Errorf("unknown http request format: %s", "foo:/bar:"),
		},
		{
			"get:/foo",
			httpRequest{
				method: "GET",
				url:    "/foo",
			},
			nil,
		},
		{
			"get:/foo:",
			httpRequest{
				method: "GET",
				url:    "/foo",
			},
			nil,
		},
		{
			"get:/foo?start={today}&end={tomorrow}:",
			httpRequest{
				method: "GET",
				url:    fmt.Sprintf("/foo?start=%s&end=%s", time.Now().Format("2006-01-02"), time.Now().Add(24*time.Hour).Format("2006-01-02")),
			},
			nil,
		},
		{
			"put:/bar?start={today-3}&end={today+4}:\\{\"date\":\"{today+20}\"\\}",
			httpRequest{
				method: "PUT",
				url:    fmt.Sprintf("/bar?start=%s&end=%s", time.Now().Add(-3*24*time.Hour).Format("2006-01-02"), time.Now().Add(4*24*time.Hour).Format("2006-01-02")),
				body:   fmt.Sprintf("\\{\"date\":\"%s\"\\}", time.Now().Add(20*24*time.Hour).Format("2006-01-02")),
			},
			nil,
		},
		{
			"put:/test:\\{\"test\":\"test\"\\}",
			httpRequest{
				method: "PUT",
				url:    "/test",
				body:   "\\{\"test\":\"test\"\\}",
			},
			nil,
		},
		{
			"post:/baz:\\{\"date\":\"{today}\"\\}",
			httpRequest{
				method: "POST",
				url:    "/baz",
				body:   fmt.Sprintf("\\{\"date\":\"%s\"\\}", time.Now().Format("2006-01-02")),
			},
			nil,
		},
		{
			"post:/test:\\{\"test\":\"test\"\\}",
			httpRequest{
				method: "POST",
				url:    "/test",
				body:   "\\{\"test\":\"test\"\\}",
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		actual, err := createHttpRequest(testCase.input)

		if (testCase.expectedErr != nil || err != nil) && (testCase.expectedErr.Error() != err.Error()) {
			t.Fatalf("Expected error\n%v\nbut got\n%v", testCase.expectedErr, err)
		} else if actual != testCase.expectedOutput {
			t.Fatalf("Expected output\n%s\nbut got\n%s", testCase.expectedOutput, actual)
		}
	}
}

func TestInterpolateDates(t *testing.T) {
	expected := fmt.Sprintf(
		"the date today is %s which is three days before %s and 2 weeks after %s and 1 day before %s",
		time.Now().Format("2006-01-02"),
		time.Now().Add(3*24*time.Hour).Format("2006-01-02"),
		time.Now().Add(-14*24*time.Hour).Format("2006-01-02"),
		time.Now().Add(24*time.Hour).Format("2006-01-02"),
	)
	actual := interpolateDates("the date today is {today} which is three days before {today+3} and 2 weeks after {today-14} and 1 day before {tomorrow}")

	if actual != expected {
		t.Fatalf("Expected \n%s\nbut got\n%s", expected, actual)
	}
}

type escapeQueryStringTestCase struct {
	input  string
	output string
}

func TestEscapeQueryStringValues(t *testing.T) {
	testCases := []escapeQueryStringTestCase{
		{
			"ab:bc:de/dodgy input that's not a [URL]",
			"ab:bc:de/dodgy input that's not a [URL]",
		},
		{
			"https://example.com/my/url",
			"https://example.com/my/url",
		},
		{
			//query params are alphabetically sorted by the golang library
			"https://example.com/my/url/page.html?surname=doe&name=john",
			"https://example.com/my/url/page.html?name=john&surname=doe",
		},
		{
			"https://example.com/my/url/page.html?item[1]=one&item[2]=two",
			"https://example.com/my/url/page.html?item%5B1%5D=one&item%5B2%5D=two",
		},
		{
			"https://example.com/my/url/page.html?request.date=2019-02-05",
			"https://example.com/my/url/page.html?request.date=2019-02-05",
		},
	}

	for _, testCase := range testCases {
		actual := escapeQueryStringValues(testCase.input)

		if actual != testCase.output {
			t.Fatalf("Expected output\n%s\nbut got\n%s\nfor url %s", testCase.output, actual, testCase.input)
		}
	}
}
