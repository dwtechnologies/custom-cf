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
	function     = "userpool-client"
	resourceType = "Custom::CognitoUserPoolClient"
	httpTimeout  = 30
)

type config struct {
	log *l.Client
	svc *cognitoidentityprovider.CognitoIdentityProvider

	physicalID            string  // The physical ID to use for the resource.
	resourceProperties    *Client // The new resource data from the template.
	oldResourceProperties *Client // The old resource data, only on updates.
}

// Client contains the data for the UserPool Client Settings.
type Client struct {
	id string

	// Standard features.
	ClientName           string                                          `json:"ClientName"` /* required */
	UserPoolID           string                                          `json:"UserPoolId"` /* required */
	GenerateSecret       bool                                            `json:"GenerateSecret"`
	RefreshTokenValidity int64                                           `json:"RefreshTokenValidity,omitempty"`
	ReadAttributes       []string                                        `json:"ReadAttributes,omitempty"`
	WriteAttributes      []string                                        `json:"WriteAttributes,omitempty"`
	ExplicitAuthFlows    []cognitoidentityprovider.ExplicitAuthFlowsType `json:"ExplicitAuthFlows,omitempty"`

	// Extended features.
	AllowedOAuthFlows               []cognitoidentityprovider.OAuthFlowType `json:"AllowedOAuthFlows,omitempty"`
	AllowedOAuthFlowsUserPoolClient bool                                    `json:"AllowedOAuthFlowsUserPoolClient,omitempty"`
	AllowedOAuthScopes              []string                                `json:"AllowedOAuthScopes,omitempty"`

	CallbackURLs       []string `json:"CallbackURLs,omitempty"`
	LogoutURLs         []string `json:"LogoutURLs,omitempty"`
	DefaultRedirectURI string   `json:"DefaultRedirectURI,omitempty"`

	SupportedIdentityProviders []string `json:"SupportedIdentityProviders,omitempty"`

	AnalyticsConfiguration *AnalyticsConfigurationType `json:"AnalyticsConfiguration,omitempty"`
}

// AnalyticsConfigurationType contains config for Analytics on the Client.
type AnalyticsConfigurationType struct {
	ApplicationID  string `json:"ApplicationId"`
	ExternalID     string `json:"ExternalId"`
	RoleArn        string `json:"RoleArn"`
	UserDataShared bool   `json:"UserDataShared"`
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

	// Set physical ID - we need the generate secret as physical id,
	// since changing this needs replacement of the resource.
	c.physicalID = fmt.Sprintf("%s-%t-%s", c.resourceProperties.UserPoolID, c.resourceProperties.GenerateSecret, c.resourceProperties.ClientName)

	// create, update or delete the userpool client.
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
		resourceProperties:    &Client{},
		oldResourceProperties: &Client{},
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

// run will either create, update or delete the specified userpool client.
// If the client already exists in the user pool it will be adopted into the
// cf stack. This so that manually created clients don't have to be recreated.
// Returns map[string]string and error.
func (c *config) run(req *events.Request) (map[string]string, error) {
	// Check for the correct ResourceType
	if req.ResourceType != resourceType {
		return nil, fmt.Errorf("Wrong ResourceType in request. Expected %s but got %s", resourceType, req.ResourceType)
	}

	// Check if the client already exists.
	client, err := c.getClientByName(c.resourceProperties.UserPoolID, c.resourceProperties.ClientName)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute API call against Cognito. Error %s", err.Error())
	}

	switch {
	// If Delete is run on the resource but the Client doesn't exist.
	case req.RequestType == "Delete" && client == nil:
		return nil, nil

	// If Delete is run on the resource.
	case req.RequestType == "Delete" && client != nil:
		return nil, c.deleteClient(req, client.id)

	// If Update is run on the resource but the Client doesn't exist
	// create it. If it was a resource that needed replacement a delete event
	// will be sent on the old resource once the new one has been created.
	case req.RequestType == "Update" && client == nil:
		return c.createClient(req)

	// If Update is run on the resource.
	case req.RequestType == "Update" && client != nil:
		// return c.updateClient(req, client.id)

	// If Create is run on the resource but the Client doesn't exist.
	case req.RequestType == "Create" && client == nil:
		return c.createClient(req)

	// If Create is run on the resource and the Client exists, adopt and update it.
	case req.RequestType == "Create" && client != nil:
		// return c.updateClient(req, client.id)
	}

	return nil, fmt.Errorf("Didn't get RequestType Create, Update or Delete")
}

// getClientByName will get the userpool client with clientName on User Pool
// with poolID. If nil is returned no client by that name was found.
// Return *Client and error.
func (c *config) getClientByName(poolID string, clientName string) (*Client, error) {
	// Validate input.
	switch {
	case poolID == "":
		return nil, fmt.Errorf("No UserPool ID specified")

	case clientName == "":
		return nil, fmt.Errorf("No Client Name specified")
	}

	// Since we need the Client ID to do any changes we first need to list
	// all clients and see if any matches our name.
	list, err := c.getClientsFromUserPool(poolID, nil, nil)
	if err != nil {
		return nil, err
	}

	// Loop over list of clients to match our clientName.
	id := ""
	for _, client := range list {
		if *client.ClientName == clientName {
			id = *client.ClientId
		}
	}

	// If clientID is empty the client doesn't exists.
	if id == "" {
		return nil, nil
	}

	resp, err := c.svc.DescribeUserPoolClientRequest(
		&cognitoidentityprovider.DescribeUserPoolClientInput{
			UserPoolId: &poolID,
			ClientId:   &id,
		}).Send()
	if err != nil {
		// If the Client doesn't exists. Return nil, nil.
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}

		return nil, err
	}

	return c.responseToClient(resp, id)
}

// getClientsFromUserPool takes poolID, clients and nextToken and retrieves all clients on
// the userpool with UserPoolID poolID. This function is recursive so it will
// execute it self if there is a nextToken. Leave nextToken as nil if it's the first run.
// Returns []cognitoidentityprovider.UserPoolClientDescription and error.
func (c *config) getClientsFromUserPool(poolID string, clients []cognitoidentityprovider.UserPoolClientDescription, nextToken *string) ([]cognitoidentityprovider.UserPoolClientDescription, error) {
	// If clients is nil, create it.
	if clients == nil {
		clients = []cognitoidentityprovider.UserPoolClientDescription{}
	}

	// Create the Input and set nextToken if it's set.
	input := &cognitoidentityprovider.ListUserPoolClientsInput{UserPoolId: &poolID}
	if nextToken != nil {
		input.NextToken = nextToken
	}

	// Get the clients for the userpool.
	resp, err := c.svc.ListUserPoolClientsRequest(input).Send()
	if err != nil {
		return clients, fmt.Errorf("Couldn't get Clients for UserPool ID: %s. Error %s", poolID, err.Error())
	}

	// Append clients.
	clients = append(clients, resp.UserPoolClients...)

	// If responses nextToken isn't nil, run recursive function.
	if resp.NextToken != nil {
		return c.getClientsFromUserPool(poolID, clients, resp.NextToken)
	}

	return clients, nil
}

// responseToClient takes resp and id and converts it to Client struct and returns it.
// Returns *Client and error.
func (c *config) responseToClient(resp *cognitoidentityprovider.DescribeUserPoolClientOutput, id string) (*Client, error) {
	// Simple validation that will result in error.
	switch {
	case resp.UserPoolClient.ClientName == nil:
		return nil, fmt.Errorf("ClientName can't be empty")

	case resp.UserPoolClient.UserPoolId == nil:
		return nil, fmt.Errorf("UserPoolId can't be empty")
	}

	client := &Client{
		id:         id,
		ClientName: *resp.UserPoolClient.ClientName,
		UserPoolID: *resp.UserPoolClient.UserPoolId,
	}

	return client, nil
}
