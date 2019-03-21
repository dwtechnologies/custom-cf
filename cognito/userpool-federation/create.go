package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// createIdentityProvider will create a new identity provider on the user pool with
// settings specified by req. If the user pool already exists it will be updated with
// the settings in req.
// Returns error.
func (c *config) createIdentityProvider(req *events.Request) (map[string]string, error) {
	switch {
	case c.resourceProperties.UserPoolID == "":
		return nil, fmt.Errorf("No UserPool ID specified")

	case c.resourceProperties.ProviderName == "":
		return nil, fmt.Errorf("No Identity Provider Name specified")
	}

	resp, err := c.svc.CreateIdentityProviderRequest(
		&cognitoidentityprovider.CreateIdentityProviderInput{
			ProviderName:     &c.resourceProperties.ProviderName,
			ProviderType:     cognitoidentityprovider.IdentityProviderTypeType(c.resourceProperties.ProviderType),
			UserPoolId:       &c.resourceProperties.UserPoolID,
			ProviderDetails:  c.resourceProperties.ProviderDetails,
			AttributeMapping: c.resourceProperties.AttributeMapping,
		}).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to create Identity Provider. Error %s", err.Error())
	}

	return map[string]string{
		"ProviderName": *resp.IdentityProvider.ProviderName,
		"ProviderType": string(resp.IdentityProvider.ProviderType),
		"UserPoolId":   *resp.IdentityProvider.UserPoolId,
	}, nil
}
