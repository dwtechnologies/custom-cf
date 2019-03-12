package events

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testRequest struct {
	req        Request
	physicalID string
	data       map[string]string
	err        error // Simulated error
	resp       string
	respErr    error // Expected error
}

type testUnmarshal struct {
	req *Request
	new *testProps
	old *testProps
}

type testProps struct {
	Key1 string `json:"Key1"`
}

// Test common Send scenarios.
func TestRespond(t *testing.T) {
	// Create response channel.
	resp := make(chan string, 4) // Update this line to correspond with the number of tests.
	// Create httptest server.
	srv := httptest.NewServer(handler(resp, t))

	// Create test cases.
	tests := []testRequest{
		// Create test request.
		testRequest{
			req: Request{
				RequestType:        "CREATE",
				StackID:            "stack1",
				ResourceType:       "Custom:TestResource",
				LogicalResourceID:  "resource1",
				RequestID:          "1234",
				ResponseURL:        srv.URL,
				ResourceProperties: []byte(`{"data1":"new"}`),
			},
			physicalID: "testId1",
			data:       map[string]string{"key1": "value1"},
			resp:       `{"Status":"SUCCESS","PhysicalResourceId":"testId1","StackId":"stack1","RequestId":"1234","LogicalResourceId":"resource1","Data":{"key1":"value1"}}`,
		},
		// Update test request.
		testRequest{
			req: Request{
				RequestType:           "UPDATE",
				StackID:               "stack2",
				ResourceType:          "Custom:TestResource",
				LogicalResourceID:     "resource2",
				RequestID:             "4321",
				ResponseURL:           srv.URL,
				ResourceProperties:    []byte(`{"data1":"updated"}`),
				OldResourceProperties: []byte(`{"data1":"old"}`),
			},
			physicalID: "testId2",
			resp:       `{"Status":"SUCCESS","PhysicalResourceId":"testId2","StackId":"stack2","RequestId":"4321","LogicalResourceId":"resource2"}`,
		},
		// Delete test request.
		testRequest{
			req: Request{
				RequestType:           "DELETE",
				StackID:               "stack3",
				ResourceType:          "Custom:TestResource",
				LogicalResourceID:     "resource3",
				RequestID:             "1111",
				ResponseURL:           srv.URL,
				OldResourceProperties: []byte(`{"data1":"old"}`),
			},
			physicalID: "testId3",
			resp:       `{"Status":"SUCCESS","PhysicalResourceId":"testId3","StackId":"stack3","RequestId":"1111","LogicalResourceId":"resource3"}`,
		},
		// Request where creation of resource FAILED.
		testRequest{
			req: Request{
				RequestType:        "CREATE",
				StackID:            "stack4",
				ResourceType:       "Custom:TestResource",
				LogicalResourceID:  "resource4",
				RequestID:          "2222",
				ResponseURL:        srv.URL,
				ResourceProperties: []byte(`{"data1":"wrong data"}`),
			},
			physicalID: "testId4",
			data:       map[string]string{"key1": "value1"},
			err:        fmt.Errorf("Couldn't create the resource. Resource with same physical id already exists"),
			resp:       `{"Status":"FAILED","Reason":"Couldn't create the resource. Resource with same physical id already exists","PhysicalResourceId":"testId4","StackId":"stack4","RequestId":"2222","LogicalResourceId":"resource4"}`,
		},
		// Request where the ResponseURL is missing.
		testRequest{
			req: Request{
				RequestType:        "CREATE",
				StackID:            "stack5",
				ResourceType:       "Custom:TestResource",
				LogicalResourceID:  "resource5",
				RequestID:          "3333",
				ResponseURL:        "",
				ResourceProperties: []byte(`{"data1":"missing url"}`),
			},
			physicalID: "testId4",
			resp:       "",
			respErr:    fmt.Errorf("Pre-signed S3 url can't be empty"),
		},
	}

	// Loop over all tests.
	for i, test := range tests {
		if err := test.req.Send(test.physicalID, test.data, test.err); err != nil {
			if err.Error() != test.respErr.Error() {
				t.Errorf("Test number: %d failed. Wanted %s but got %s", i+1, test.respErr, err.Error())
			}
			continue
		}

		// Check that we got the same response.
		val := <-resp
		if val != test.resp {
			t.Errorf("Test number: %d failed. Wanted %s but got %s", i+1, test.resp, val)
		}
	}
}

// Handler for sending 200's and sending the incoming request body to resp channel.
func handler(resp chan string, t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		defer r.Body.Close()
		resp <- string(body)
	})
}

// Handler for sending 500 error.
func handlerError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})
}

// Test that sending sendResponse with an empty body and request
// against malformed URL fails.
func TestSendResponse(t *testing.T) {
	req := &Request{
		ResponseURL: `h\t\t\p:\//wron.com?this=wrong?=this`,
	}

	// Test with empty body.
	if err := req.sendResponse(nil, 500); err != nil {
		if err.Error() != "Body of response can't be empty" {
			t.Errorf("%s", err.Error())
		}
	}

	// Test the bogus url.
	if err := req.sendResponse([]byte("{}"), 500); err != nil {
		if err.Error() != `Couldn't create request for s3-presigned-url. Error parse h\t\t\p:\//wron.com?this=wrong?=this: first path segment in URL cannot contain colon` {
			t.Errorf("%s", err.Error())
		}
	}
}

// Test that request to the wrong URL is handled correctly.
func TestDoRequestWrongUrl(t *testing.T) {
	req := &Request{}
	saveReq, err := http.NewRequest("PUT", "http://127.0.0.1:1/wrong/url", nil)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	// Check that error works.
	if err := req.doRequest(&http.Client{Timeout: time.Duration(10) * time.Millisecond}, saveReq, 0); err != nil {
		errMessage := "Couldn't send request for s3-presigned-url. Attempt 5. Error Put http://127.0.0.1:1/wrong/url: dial tcp 127.0.0.1:1: connect: connection refused"
		if err.Error() != errMessage {
			t.Errorf("Expected '%s' but got '%s'", errMessage, err.Error())
		}
	}
}

// Test that we retry the request for 5 times before giving up.
func TestDoRequestTestRetries(t *testing.T) {
	srv := httptest.NewServer(handlerError())

	req := &Request{}
	saveReq, err := http.NewRequest("PUT", srv.URL, nil)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	// Check that retry works and that it fails after 5 attempts.
	if err := req.doRequest(&http.Client{Timeout: time.Duration(10) * time.Millisecond}, saveReq, 0); err != nil {
		errMessage := "Didn't receive error, but response wasn't 200. Status Code: 500. Attempt 5"
		if err.Error() != errMessage {
			t.Errorf("Expected '%s' but got '%s'", errMessage, err.Error())
		}
	}
}

// Will test Unmarshal.
func TestUnmarshal(t *testing.T) {
	req := &Request{RequestType: "Update", ResourceProperties: []byte(`{"key1":"value1"}`), OldResourceProperties: []byte(`{"key1":"value2"}`)}
	new := &testProps{}
	old := &testProps{}

	if err := req.Unmarshal(new, old); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	if new.Key1 != "value1" {
		t.Errorf("Expected value1 but got %s", new.Key1)
	}

	if old.Key1 != "value2" {
		t.Errorf("Expected value2 but got %s", old.Key1)
	}
}

// Will test Unmarshal when not Type update. OldResourceProperties should always be nil.
func TestUnmarshalNotUpdate(t *testing.T) {
	req := &Request{RequestType: "Create", ResourceProperties: []byte(`{"key1":"value1"}`), OldResourceProperties: []byte(`{"key1":"value2"}`)}
	new := &testProps{}
	old := &testProps{}

	if err := req.Unmarshal(new, old); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	if new.Key1 != "value1" {
		t.Errorf("Expected value1 but got %s", new.Key1)
	}

	if old.Key1 != "" {
		t.Errorf("Expected empty for old value")
	}
}

// Test Unmarshal when no values are set.
func TestUnmarshalNilValues(t *testing.T) {
	req := &Request{}
	new := &testProps{}
	old := &testProps{}

	if err := req.Unmarshal(new, old); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	if new.Key1 != "" {
		t.Errorf("Expected empty value but got %s", new.Key1)
	}

	if old.Key1 != "" {
		t.Errorf("Expected empty value but got %s", old.Key1)
	}
}

// Test Unmarshal when values are set but empty.
func TestUnmarshalEmptyValues(t *testing.T) {
	req := &Request{ResourceProperties: []byte(""), OldResourceProperties: []byte("")}
	new := &testProps{}
	old := &testProps{}

	if err := req.Unmarshal(new, old); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	if new.Key1 != "" {
		t.Errorf("Expected empty value but got %s", new.Key1)
	}

	if old.Key1 != "" {
		t.Errorf("Expected empty value but got %s", old.Key1)
	}
}

// Unmarshal wrong values into new.
func TestUnmarshalWrongValuesNew(t *testing.T) {
	req := &Request{ResourceProperties: []byte(`{"Key1":123}`)}
	new := &testProps{}

	if err := req.Unmarshal(new, nil); err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

// Unmarshal wrong values into old.
func TestUnmarshalWrongValuesOld(t *testing.T) {
	req := &Request{RequestType: "Update", OldResourceProperties: []byte(`{"Key1":123}`)}
	old := &testProps{}

	if err := req.Unmarshal(nil, old); err == nil {
		t.Errorf("Expected error, but got nil")
	}
}
