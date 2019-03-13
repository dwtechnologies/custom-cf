package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// createClient will create a new client on the user pool with settings specified by req.
// Returns map of string that is data that Fn::GetAtt can use.
// Returns map[string]string error.
func (c *config) createClient(req *events.Request) (map[string]string, error) {
	input := &cognitoidentityprovider.CreateUserPoolClientInput{
		ClientName:                      &c.resourceProperties.ClientName,
		UserPoolId:                      &c.resourceProperties.UserPoolID,
		GenerateSecret:                  &c.resourceProperties.GenerateSecret,
		RefreshTokenValidity:            c.resourceProperties.RefreshTokenValidity,
		ReadAttributes:                  c.resourceProperties.ReadAttributes,
		WriteAttributes:                 c.resourceProperties.WriteAttributes,
		ExplicitAuthFlows:               c.resourceProperties.ExplicitAuthFlows,
		AllowedOAuthFlows:               c.resourceProperties.AllowedOAuthFlows,
		AllowedOAuthFlowsUserPoolClient: c.resourceProperties.AllowedOAuthFlowsUserPoolClient,
		AllowedOAuthScopes:              c.resourceProperties.AllowedOAuthScopes,
		CallbackURLs:                    c.resourceProperties.CallbackURLs,
		LogoutURLs:                      c.resourceProperties.LogoutURLs,
		DefaultRedirectURI:              c.resourceProperties.DefaultRedirectURI,
		SupportedIdentityProviders:      c.resourceProperties.SupportedIdentityProviders,
	}

	// Add analytics if it's not nil.
	if c.resourceProperties.AnalyticsConfiguration != nil {
		input.AnalyticsConfiguration = &cognitoidentityprovider.AnalyticsConfigurationType{
			ApplicationId:  c.resourceProperties.AnalyticsConfiguration.ApplicationID,
			ExternalId:     c.resourceProperties.AnalyticsConfiguration.ExternalID,
			RoleArn:        c.resourceProperties.AnalyticsConfiguration.RoleArn,
			UserDataShared: c.resourceProperties.AnalyticsConfiguration.UserDataShared,
		}
	}

	resp, err := c.svc.CreateUserPoolClientRequest(input).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to create Client. Error %s", err.Error())
	}

	return map[string]string{
		"ClientName":   *resp.UserPoolClient.ClientName,
		"ClientId":     *resp.UserPoolClient.ClientId,
		"ClientSecret": *resp.UserPoolClient.ClientSecret,
		"UserPoolId":   *resp.UserPoolClient.UserPoolId,
	}, nil
}
