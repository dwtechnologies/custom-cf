// Package respond can be used to handle requests for custom resources for
// CloudFormation as well as handling the response and saving the results
// of the operation to the pre-signed S3 URL.
package respond

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Request from CloudFormation.
// You will need to Unmarshal ResourceProperties and OldResourceProperties
// into your own struct to be usable.
type Request struct {
	RequestType           string          `json:"RequestType"`
	ResponseURL           string          `json:"ResponseURL"`
	StackID               string          `json:"StackId"`
	RequestID             string          `json:"RequestId"`
	ResourceType          string          `json:"ResourceType"`
	LogicalResourceID     string          `json:"LogicalResourceId"`
	ResourceProperties    json.RawMessage `json:"ResourceProperties"`
	OldResourceProperties json.RawMessage `json:"OldResourceProperties"`
}

// Response is the data that will be stored on the pre-signed S3 url.
// The Data part should contain your ResourceProperties as JSON
type response struct {
	Status             string            `json:"Status"`             /* Required */
	Reason             string            `json:"Reason,omitempty"`   /* Should only be set if Status == Failed */
	PhysicalResourceID string            `json:"PhysicalResourceId"` /* Required */
	StackID            string            `json:"StackId"`            /* Required */
	RequestID          string            `json:"RequestId"`          /* Required */
	LogicalResourceID  string            `json:"LogicalResourceId"`  /* Required */
	Data               map[string]string `json:"Data,omitempty"`     /* Resource Properties data that can be accessed through Fn::GatAtt*/
}

// Send takes physicalID and respError and sends it to an S3 pre-signed url.
// physicalID should be a unique physicalID that the resource should have, naming will
// depend on the type of resource you're creating but can often be "put together" by
// various fields from ResourceProperties.
// data should be key value pairs of data that you want to be able to access with the
// Fn::GetAtt function in the CloudFormation template.
// respErr is the response error, if the resource creation failed we still need to save
// the state FAILED to S3 for the Custom Resource to work.
// Returns error.
func (req *Request) Send(physicalID string, data map[string]string, respErr error) error {
	// Create Response.
	body, err := req.createResponse(physicalID, data, respErr)
	if err != nil {
		return err
	}

	// Send the response.
	if err := req.sendResponse(body, 30000); err != nil {
		return err
	}

	return nil
}

// createResponse takes physicalID, data and err and creates a response JSON bytes that can be sent to the pre-signed s3 url.
// Where data is key value pairs that you want to be accessed through Fn::GetAtt, can be nil if not needed.
// Returns []byte and error.
func (req *Request) createResponse(physicalID string, data map[string]string, respErr error) ([]byte, error) {
	resp := &response{
		Status:             "SUCCESS",
		StackID:            req.StackID,
		RequestID:          req.RequestID,
		PhysicalResourceID: physicalID,
		LogicalResourceID:  req.LogicalResourceID,
		Data:               data,
	}

	if respErr != nil {
		resp.Status = "FAILED"
		resp.Reason = respErr.Error()
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("Couldn't JSON Marshal the Response. Error %s", err.Error())
	}

	return b, nil
}

// sendResponse created a response for the s3-presigned-url and sends it.
// It sets the http client timeout to timeOut in milliseconds.
// Returns error.
func (req *Request) sendResponse(body []byte, timeOut int) error {
	client := &http.Client{Timeout: time.Duration(timeOut) * time.Millisecond}

	switch {
	case req.ResponseURL == "":
		return fmt.Errorf("Pre-signed S3 url can't be empty")

	case len(body) == 0:
		return fmt.Errorf("Body of response can't be empty")
	}

	// Create the request for s3-presigned-url.
	saveReq, err := http.NewRequest("PUT", req.ResponseURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Couldn't create request for s3-presigned-url. Error %s", err.Error())
	}

	// Send request for s3-presigned-url.
	return req.doRequest(client, saveReq, 0)
}

// doRequest takes http.Client, http.Request and attempt number (should be 0 for first)
// and will run as a recursive function for 5 attempts if it fails.
func (req *Request) doRequest(client *http.Client, saveReq *http.Request, attempt int) error {
	resp, err := client.Do(saveReq)

	// Set the correct error message.
	switch {
	case err != nil:
		err = fmt.Errorf("Couldn't send request for s3-presigned-url. Attempt %d. Error %s", attempt+1, err.Error())

	case err == nil && resp.StatusCode != 200:
		err = fmt.Errorf("Didn't receive error, but response wasn't 200. Status Code: %d. Attempt %d", resp.StatusCode, attempt+1)
	}

	// If error is nit nil and we haven't retried 5 times test again.
	// Otherwise return error.
	switch {
	case err != nil && attempt < 4:
		return req.doRequest(client, saveReq, attempt+1)

	case err != nil:
		return err
	}

	return nil
}
