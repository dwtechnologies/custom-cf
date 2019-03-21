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
	function     = "userpool-federation"
	resourceType = "Custom::CognitoUserPoolFederation"
	httpTimeout  = 30
)

type config struct {
	log *l.Client
	svc *cognitoidentityprovider.CognitoIdentityProvider

	physicalID string // The physical ID to use for the resource.
	resourceProperties    *IdentityProvider // The new resource data from the template.
	oldResourceProperties *IdentityProvider // The old resource data, only on updates.
}

// IdentityProvider valid ProviderTypes are
// SAML, Facebook, Google, LoginWithAmazon or OIDC
type IdentityProvider struct {
	IdpIdentifiers   []string          `json:"-"`
	ProviderName     string            `json:"ProviderName"`
	ProviderType     string            `json:"ProviderType"`
	ProviderDetails  map[string]string `json:"ProviderDetails"`
	AttributeMapping map[string]string `json:"AttributeMapping"`
	UserPoolID       string            `json:"UserPoolId"`
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
	c.physicalID = fmt.Sprintf("%s-%s", c.resourceProperties.UserPoolID, c.resourceProperties.ProviderName)

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
		}),
		physicalID: "NotAviable",
		resourceProperties: &IdentityProvider{},
		oldResourceProperties: &IdentityProvider{},
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

// run will either create, update or delete the specified identity provider.
// If the identity provider already exists in the user pool it will be adopted into the
// cf stack. This so that manually created identity providers don't have to be recreated.
// Returns map[string]string and error.
func (c *config) run(req *events.Request) (map[string]string, error) {
	// Check for the correct ResourceType
	if req.ResourceType != resourceType {
		return nil, fmt.Errorf("Wrong ResourceType in request. Expected %s but got %s", resourceType, req.ResourceType)
	}

	// Check if the Identity Provider already exists.
	provider, err := c.getIdentityProviderByName(c.resourceProperties.UserPoolID, c.resourceProperties.ProviderName)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute API call against Cognito. Error %s", err.Error())
	}

	switch {
	// If Delete is run on the stack but the Identity Provider doesn't exist.
	case req.RequestType == "Delete" && provider == nil:
		return nil, nil

	// If Delete is run on the stack.
	case req.RequestType == "Delete" && provider != nil:
		return nil, c.deleteIdentityProvider(req)

	// If Update is run on the stack but the Identity Provider doesn't exist.
	// Create it. If it was a resource that needed replacement a delete event
	// will be sent on the old resource once the new one has been created.
	case req.RequestType == "Update" && provider == nil:
		return c.createIdentityProvider(req)

	// If Update is run on the stack.
	case req.RequestType == "Update" && provider != nil:
		return c.updateIdentityProvider(req, provider.IdpIdentifiers)

	// If Create is run on the stack but the Identity Provider doesn't exist.
	case req.RequestType == "Create" && provider == nil:
		return c.createIdentityProvider(req)

	// If Create is run on the stack and the Identity Provider exists, adopt and update it.
	case req.RequestType == "Create" && provider != nil:
		return c.updateIdentityProvider(req, provider.IdpIdentifiers)
	}

	return nil, fmt.Errorf("Didn't get RequestType Create, Update or Delete")
}

// getIdentityProviderByName will get the identity provider with providerName on User Pool
// with poolID. If nil is returned no identity provider by that name was found.
// Return *IdentityProvider and error.
func (c *config) getIdentityProviderByName(poolID string, providerName string) (*IdentityProvider, error) {
	// Just return nil, nil if any of the required fields are missing.
	// Extra validation will be done in the specific resource creation
	// functions. This is so that Delete on empty will not fail.
	switch {
	case poolID == "":
		return nil, nil

	case providerName == "":
		return nil, nil
	}

	resp, err := c.svc.DescribeIdentityProviderRequest(
		&cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   &poolID,
			ProviderName: &providerName,
		}).Send()
	if err != nil {
		// If the Identity Provier doesn't exists. Return nil and no error.
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}

		return nil, err
	}

	return &IdentityProvider{
		IdpIdentifiers:   resp.IdentityProvider.IdpIdentifiers,
		ProviderName:     *resp.IdentityProvider.ProviderName,
		ProviderType:     string(resp.IdentityProvider.ProviderType),
		ProviderDetails:  resp.IdentityProvider.ProviderDetails,
		UserPoolID:       *resp.IdentityProvider.UserPoolId,
		AttributeMapping: resp.IdentityProvider.AttributeMapping,
	}, nil
}
