package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/dwtechnologies/custom-cf/lib/events"
	l "github.com/nuttmeister/llogger"
)

// http client timeout in seconds.
const (
	service      = "custom-cf"
	function     = "tag"
	resourceType = "Custom::IAMRoleTags"
	httpTimeout  = 30
)

type config struct {
	log *l.Client
	svc *iam.IAM

	physicalID            string    // The physical ID to use for the resource.
	resourceProperties    *RoleTags // The new resource data from the template.
	oldResourceProperties *RoleTags // The old resource data, only on updates.
}

type RoleTags struct {
	RoleName string    `json:"RoleName"`
	Tags     []iam.Tag `json:"Tags"`
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

	// Create the Cognit service.
	if err := c.createCognitoService(); err != nil {
		return c.runError(req, err)
	}

	// Unmarshal the ResourceProperties and OldResourceProperties into the config.
	if err := req.Unmarshal(c.resourceProperties, c.oldResourceProperties); err != nil {
		return c.runError(req, err)
	}

	// Set physical ID
	c.physicalID = fmt.Sprintf("%s-tag", c.resourceProperties.RoleName)

	// create, update or delete
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

// createConfig takes ctx and req and creates a config that contains the logger
// and the cognito aws service.
// Returns *config and error.
func createConfig(ctx context.Context, req *events.Request) *config {
	return &config{
		log: l.Create(ctx, l.Input{
			"service":               service,
			"function":              function,
			"env":                   os.Getenv("ENVIRONMENT"),
			"stackId":               req.StackID,
			"requestType":           req.RequestType,
			"requestId":             req.RequestID,
			"resourceType":          req.ResourceType,
			"logicalResourceId":     req.LogicalResourceID,
			"resourceProperties":    req.ResourceProperties,
			"oldResourceProperties": req.OldResourceProperties,
		}),
		physicalID:            "NotAviable",
		resourceProperties:    &RoleTags{},
		oldResourceProperties: &RoleTags{},
	}
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

// Creates AWS Config and the CognitoIdentity Service.
// Returns error.
func (c *config) createCognitoService() error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return fmt.Errorf("Couldn't create AWS cfg. Error: %s", err.Error())
	}

	c.svc = iam.New(cfg)
	return nil
}

// run will either create, update or delete the UI customization
// Returns map[string]string and error.
func (c *config) run(req *events.Request) (map[string]string, error) {
	// Check for the correct ResourceType
	if req.ResourceType != resourceType {
		return nil, fmt.Errorf("Wrong ResourceType in request. Expected %s but got %s", resourceType, req.ResourceType)
	}

	switch {
	case req.RequestType == "Delete":
		return nil, c.deleteTags(req)

	case req.RequestType == "Create":
		return c.createTags(req)

	case req.RequestType == "Update":
		return c.updateTags(req)
	}

	return nil, fmt.Errorf("Didn't get RequestType Create, Update or Delete")
}
