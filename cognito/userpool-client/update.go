package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// updateClient will update the client on the user pool with settings specified by req.
// Returns map of string that is data that Fn::GetAtt can use.
// Returns map[string]string error.
func (c *config) updateClient(req *events.Request, id string) (map[string]string, error) {

	// Simple validation that will result in error.
	switch {
	case id == "":
		return nil, fmt.Errorf("ClientId can't be empty")

	case c.resourceProperties.ClientName == "":
		return nil, fmt.Errorf("ClientName can't be empty")

	case c.resourceProperties.UserPoolID == "":
		return nil, fmt.Errorf("UserPoolId can't be empty")
	}

	// Set oauthFlows
	oauthFlows := false
	if c.resourceProperties.AllowedOAuthFlowsUserPoolClient == "true" {
		oauthFlows = true
	}

	input := &cognitoidentityprovider.UpdateUserPoolClientInput{
		ClientId:                        &id,
		ClientName:                      &c.resourceProperties.ClientName,
		UserPoolId:                      &c.resourceProperties.UserPoolID,
		AllowedOAuthFlowsUserPoolClient: &oauthFlows,
	}

	// Set optional settings.
	if c.resourceProperties.RefreshTokenValidity != "" {
		n, err := strconv.ParseInt(c.resourceProperties.RefreshTokenValidity, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RefreshTokenValidity invalid")
		}
		input.RefreshTokenValidity = &n
	}
	if c.resourceProperties.ReadAttributes != nil {
		input.ReadAttributes = c.resourceProperties.ReadAttributes
	}
	if c.resourceProperties.WriteAttributes != nil {
		input.WriteAttributes = c.resourceProperties.WriteAttributes
	}
	if c.resourceProperties.ExplicitAuthFlows != nil {
		input.ExplicitAuthFlows = c.resourceProperties.ExplicitAuthFlows
	}
	if c.resourceProperties.AllowedOAuthFlows != nil {
		input.AllowedOAuthFlows = c.resourceProperties.AllowedOAuthFlows
	}
	if c.resourceProperties.AllowedOAuthScopes != nil {
		input.AllowedOAuthScopes = c.resourceProperties.AllowedOAuthScopes
	}
	if c.resourceProperties.CallbackURLs != nil {
		input.CallbackURLs = c.resourceProperties.CallbackURLs
	}
	if c.resourceProperties.LogoutURLs != nil {
		input.LogoutURLs = c.resourceProperties.LogoutURLs
	}
	if c.resourceProperties.DefaultRedirectURI != "" {
		input.DefaultRedirectURI = &c.resourceProperties.DefaultRedirectURI
	}
	if c.resourceProperties.SupportedIdentityProviders != nil {
		input.SupportedIdentityProviders = c.resourceProperties.SupportedIdentityProviders
	}

	// Set Analytics
	if c.resourceProperties.AnalyticsConfiguration != nil {
		input.AnalyticsConfiguration = &cognitoidentityprovider.AnalyticsConfigurationType{UserDataShared: &c.resourceProperties.AnalyticsConfiguration.UserDataShared}

		if c.resourceProperties.AnalyticsConfiguration.ApplicationID != "" {
			input.AnalyticsConfiguration.ApplicationId = &c.resourceProperties.AnalyticsConfiguration.ApplicationID
		}
		if c.resourceProperties.AnalyticsConfiguration.ExternalID != "" {
			input.AnalyticsConfiguration.ExternalId = &c.resourceProperties.AnalyticsConfiguration.ExternalID
		}
		if c.resourceProperties.AnalyticsConfiguration.RoleArn != "" {
			input.AnalyticsConfiguration.RoleArn = &c.resourceProperties.AnalyticsConfiguration.RoleArn
		}
	}

	resp, err := c.svc.UpdateUserPoolClientRequest(input).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to update Client. Error %s", err.Error())
	}

	attr := map[string]string{
		"ClientName": *resp.UserPoolClient.ClientName,
		"ClientId":   *resp.UserPoolClient.ClientId,
		"UserPoolId": *resp.UserPoolClient.UserPoolId,
	}

	// Only set ClientSecret if it's set.
	if resp.UserPoolClient.ClientSecret != nil {
		attr["ClientSecret"] = *resp.UserPoolClient.ClientSecret
	}

	return attr, nil
}
