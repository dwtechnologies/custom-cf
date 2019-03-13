package main

import (
	"context"
	"fmt"
	"os"

	// External
	"github.com/dwtechnologies/custom-cf/lib/events"
	l "github.com/nuttmeister/llogger"

	// External - AWS
	"github.com/aws/aws-lambda-go/lambda"
)

// http client timeout in seconds.
const (
	service      = "custom-cf"
	function     = "myresource"         // Replace with the function name
	resourceType = "Custom::MyResource" // Change to the Resource Name you want to use.
	httpTimeout  = 30
)

type config struct {
	log *l.Client
	svc interface{} // Replace with AWS service (or other service etc) that the resource needs access to.

	physicalID            string              // The physical ID to use for the resource.
	resourceProperties    *ResourceProperties // The new resource data from the template.
	oldResourceProperties *ResourceProperties // The old resource data, only on updates.
}

// ResourceProperties needs to be exported so that the lib/events package can Unmarshal it.
// This should contain all the fields that you can add in the Custom resource in CloudFormation.
type ResourceProperties struct {
	MyResourceField1 string `json:"MyResourceField1"`
	MyResourceField2 string `json:"MyResourceField2"`
}

func main() {
	lambda.Start(handler)
}

// handler takes context.Context and *events.Request.
// Returns error.
func handler(ctx context.Context, req *events.Request) error {
	// Create the config.
	c := createConfig(ctx, req)
	c.log.Print(l.Input{"loglevel": "info", "message": "Function started"})
	defer c.log.Print(l.Input{"loglevel": "info", "message": "Function finished"})

	// Function to create the AWS service (if needed) and set it to c.svc.
	// If creation of the service fails it should return c.runError(req, err).

	// Unmarshal the ResourceProperties and OldResourceProperties into the config.
	// If it fails send FAILED status back to CF and return error.
	if err := req.Unmarshal(c.resourceProperties, c.oldResourceProperties); err != nil {
		return c.runError(req, err)
	}

	// Set physical ID for the resource. Please note that if the physical ID differs
	// when an update is being run a new resource will be created and a delete event
	// will be run on the previous physical id.
	// This is how you control when a resource needs replacement instead of just pure
	// updating it.
	c.physicalID = fmt.Sprintf("%s-%s", c.resourceProperties.MyResourceField1, c.resourceProperties.MyResourceField2)

	// create, update or delete the resource.
	data, err := c.run(req)
	if err != nil {
		return c.runError(req, err)
	}

	// Send the result to the pre-signed s3 url.
	if err := req.Send(c.physicalID, data, err); err != nil {
		return err
	}
	return nil
}

// runError takes error and logs it and sends a failure request to the s3 pre-signed url.
// Returns error.
func (c *config) runError(req *events.Request, err error) error {
	if err != nil {
		c.log.Print(l.Input{"loglevel": "error", "message": err.Error()})
		if err := req.Send(c.physicalID, nil, err); err != nil {
			c.log.Print(l.Input{"loglevel": "error", "message": err.Error()})
		}
	}
	return err
}

// createConfig takes ctx and req and creates a config that contains the logger
// and the cognito aws service.
// Returns *config and error.
func createConfig(ctx context.Context, req *events.Request) *config {
	return &config{
		log: l.Create(ctx, l.Input{
			"service":            service,
			"function":           function,
			"env":                os.Getenv("ENVIRONMENT"),
			"stackId":            req.StackID,
			"requestType":        req.RequestType,
			"requestId":          req.RequestID,
			"resourceType":       req.ResourceType,
			"logicalResourceId":  req.LogicalResourceID,
			"resourceProperties": req.ResourceProperties,
		}),
		physicalID:            "NotAviable",
		resourceProperties:    &ResourceProperties{},
		oldResourceProperties: &ResourceProperties{},
	}
}

// run will either create, update or delete the specified resource.
// the map[string]string that is returned is key values on what can
// be obtained by CloudFormation Fn::GetAtt function, so that other
// resources can reference data from this resource.
// Returns map[string]string and error.
func (c *config) run(req *events.Request) (map[string]string, error) {
	// Check for the correct ResourceType.
	if req.ResourceType != resourceType {
		return nil, fmt.Errorf("Wrong ResourceType in request. Expected %s but got %s", resourceType, req.ResourceType)
	}

	// Add logic for checking if the resource with the same data already exists.
	// This is just a placeholder variable.
	exists := false

	switch {
	// If Delete is run on the stack but the resource doesn't exist / already deleted.
	case req.RequestType == "Delete" && !exists:
		return nil, nil

	// If Delete is run on the stack.
	case req.RequestType == "Delete" && exists:
		// Add logic to delete resource here.
		err := fmt.Errorf("placeholder result")
		return nil, err

	// If Update is run on the stack but the resource doesn't exists
	// create it. If it was a resource that needed replacement a delete event
	// will be sent on the old resource once the new one has been created.
	case req.RequestType == "Update" && !exists:
		// Add logic to create resource here.
		err := fmt.Errorf("placeholder result")
		return map[string]string{"key1": "value1"}, err

	// If Update is run on the stack.
	case req.RequestType == "Update" && exists:
		// Add logic to update resource here.
		err := fmt.Errorf("placeholder result")
		return map[string]string{"key1": "value1"}, err

	// If Create is run on the stack and the resource doesn't exist:
	case req.RequestType == "Create" && !exists:
		// Add logic to create resource here.
		err := fmt.Errorf("placeholder result")
		return map[string]string{"key1": "value1"}, err

	// If Create is run on the stack and the resource exists, adopt and update it.
	case req.RequestType == "Create" && exists:
		// Add logic to update resource here.
		err := fmt.Errorf("placeholder result")
		return map[string]string{"key1": "value1"}, err
	}

	return nil, fmt.Errorf("Didn't get RequestType Create, Update or Delete")
}
