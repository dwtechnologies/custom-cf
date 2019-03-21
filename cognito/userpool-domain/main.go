package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	function     = "userpool-domain"
	resourceType = "Custom::CognitoUserPoolDomain"
	httpTimeout  = 30
)

type config struct {
	log *l.Client
	svc *cognitoidentityprovider.CognitoIdentityProvider

	physicalID            string  // The physical ID to use for the resource.
	resourceProperties    *Domain // The new resource data from the template.
	oldResourceProperties *Domain // The old resource data, only on updates.
}

// Domain contains the fields for creating a UserPool Domain.
type Domain struct {
	Domain             string              `json:"Domain"`
	CustomDomainConfig *CustomDomainConfig `json:"CustomDomainConfig,omitempty"`
	UserPoolID         string              `json:"UserPoolId"`
}

// CustomDomainConfig contains the custom domain configuration.
type CustomDomainConfig struct {
	CertificateArn string `json:"CertificateArn"`
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
	c.physicalID = fmt.Sprintf("%s", c.resourceProperties.Domain)

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
		resourceProperties:    &Domain{},
		oldResourceProperties: &Domain{},
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

// run will either create, update or delete the specified domain.
// If the domain already exists in the user pool it will be adopted into the
// cf stack. This so that manually created domains don't have to be recreated.
// Returns map[string]string and error.
func (c *config) run(req *events.Request) (map[string]string, error) {
	// Check for the correct ResourceType
	if req.ResourceType != resourceType {
		return nil, fmt.Errorf("Wrong ResourceType in request. Expected %s but got %s", resourceType, req.ResourceType)
	}

	// Check if the Domain already exists.
	domain, err := c.getDomain(false)
	if err != nil {
		return nil, err
	}

	switch {
	// If Delete is run on the stack but the domain doesn't exist.
	case req.RequestType == "Delete" && domain == nil:
		return nil, nil

	// If Delete is run on the stack.
	case req.RequestType == "Delete" && domain != nil:
		return nil, c.deleteDomain(req, false)

	// If Update is run on the stack but the domain doesn't exist
	// create it. If it was a resource that needed replacement a delete event
	// will be sent on the old resource once the new one has been created.
	case req.RequestType == "Update" && domain == nil:
		oldDomain, err := c.getDomain(true)
		if err != nil {
			return nil, err
		}

		// If oldDomain is nil, the old domain has already been deleted.
		// So just create the new one.
		if oldDomain == nil {
			return c.createDomain(req)
		}
		// Update the domain.
		return c.updateDomain(req)

	// If Update is run on the stack.
	case req.RequestType == "Update" && domain != nil:
		_, err := c.getDomain(true)
		if err != nil {
			return nil, err
		}
		return c.updateDomain(req)

	// If Create is run on the stack but the domain doesn't exist.
	case req.RequestType == "Create" && domain == nil:
		return c.createDomain(req)

	// If Create is run on the stack and the domain exists, adopt and update it.
	case req.RequestType == "Create" && domain != nil:
		return c.updateDomain(req)
	}

	return nil, fmt.Errorf("Didn't get RequestType Create, Update or Delete")
}

// getDomain will get the domain with the domain specified i c.resourceProperties.Domain
// or c.oldResourceProperties.Domain depending on if old is true or false.
// If nil is returned no domain by that name was found.
// Return *Domain and error.
func (c *config) getDomain(old bool) (*Domain, error) {
	props := c.resourceProperties
	if old {
		props = c.oldResourceProperties
	}

	// Just return nil, nil if any of the required fields are missing.
	// Extra validation will be done in the specific resource creation
	// functions. This is so that Delete on empty will not fail.
	switch {
	case props.Domain == "":
		return nil, nil
	}

	resp, err := c.svc.DescribeUserPoolDomainRequest(
		&cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: &props.Domain,
		}).Send()
	if err != nil {
		// If the domain doesn't exists. Return nil and no error.
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}

		return nil, err
	}

	if resp != nil {
		fmt.Printf("%+v", *resp)
	}

	// If domain is nil, the domain doesn't exists.
	if resp.DomainDescription.Domain == nil {
		return nil, nil
	}

	// Check that the domain belongs to our UserPoolID.
	if *resp.DomainDescription.UserPoolId != props.UserPoolID {
		return nil, fmt.Errorf("Domain name exists but doesn't belong to UserPoolId: %s. But belongs to UserPoolId: %s", props.UserPoolID, *resp.DomainDescription.UserPoolId)
	}

	domain := &Domain{
		UserPoolID: *resp.DomainDescription.UserPoolId,
	}

	// Only set CustomDomainConfig if it's not nil.
	if resp.DomainDescription.CustomDomainConfig.CertificateArn != nil {
		domain.CustomDomainConfig = &CustomDomainConfig{CertificateArn: *resp.DomainDescription.CustomDomainConfig.CertificateArn}
	}

	return domain, nil
}
