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
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// http client timeout in seconds.
const (
	service      = "custom-cf"
	function     = "userpool-mfa"
	resourceType = "Custom::CognitoUserPoolMFA"
	httpTimeout  = 30
)

type config struct {
	log *l.Client
	svc *cognitoidentityprovider.CognitoIdentityProvider

	physicalID string // The physical ID to use for the resource.
	resourceProperties    *MFA // The new resource data from the template.
	oldResourceProperties *MFA // The old resource data, only on updates.
}

// MFA contains the fields for setting a UserPools MFA settings.
type MFA struct {
	MfaConfiguration     string            `json:"MfaConfiguration"`
	SmsMfaConfiguration *SmsMfaConfiguration `json:"SmsMfaConfiguration"`
	SoftwareTokenMfaConfiguration *SoftwareTokenMfaConfiguration `json:"SoftwareTokenMfaConfiguration"`
	UserPoolID       string            `json:"UserPoolId"`
}

// SmsMfaConfiguration contains the SMS MFA configuration.
type SmsMfaConfiguration struct {
	SmsAuthenticationMessage string `json:"SmsAuthenticationMessage"`
	SmsConfiguration *SmsConfiguration `json:"SmsConfiguration"`
}

// SmsConfiguration contains the configuration for sending SMS.
type SmsConfiguration struct {
	SnsCallerArn string `json:"SnsCallerArn"`
		ExternalID string `json:"ExternalId"`
}

// SoftwareTokenMfaConfiguration contains the Software MFA configuration.
type SoftwareTokenMfaConfiguration struct {
	Enabled bool `json:"Enabled"`
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
	c.physicalID = fmt.Sprintf("%s-mfa", c.resourceProperties.UserPoolID)

	// create, update or delete the userpool federation.
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
			"oldResourceProperties": req.OldResourceProperties,
		}),
		physicalID: "NotAviable",
		resourceProperties: &MFA{},
		oldResourceProperties: &MFA{},
	}
}

// Creates AWS Config and the CognitoIdentity Service.
// Returns error.
func (c *config) createCognitoService() error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return fmt.Errorf("Couldn't create AWS cfg. Error: %s", err.Error())
	}

	c.svc = cognitoidentityprovider.New(cfg)
	return nil
}

// run will either set or delete the specified MFA settings on the UserPool.
// If the domain already exists in the user pool it will be adopted into the
// cf stack. This so that manually created domains don't have to be recreated.
// Returns map[string]string and error.
func (c *config) run(req *events.Request) (map[string]string, error) {
	// Check for the correct ResourceType
	if req.ResourceType != resourceType {
		return nil, fmt.Errorf("Wrong ResourceType in request. Expected %s but got %s", resourceType, req.ResourceType)
	}

	switch {
	// If Delete is run on the stack.
	case req.RequestType == "Delete":
		return nil, c.deleteMFA(req)

	// If Update is run on the stack.
	case req.RequestType == "Update":
		return c.setMFA(req)

	// If Create is run on the stack.
	case req.RequestType == "Create":
		return c.setMFA(req)
	}

	return nil, fmt.Errorf("Didn't get RequestType Create, Update or Delete")
}
