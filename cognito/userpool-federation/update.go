package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// updateIdentityProvider updates the identity provider specified by IdpIdentifiers in provider.
// Returns error.
func (c *config) updateIdentityProvider(req *events.Request, idps []string) (map[string]string, error) {
	resp, err := c.svc.UpdateIdentityProviderRequest(
		&cognitoidentityprovider.UpdateIdentityProviderInput{
			IdpIdentifiers:   idps,
			ProviderName:     &c.resourceProperties.ProviderName,
			UserPoolId:       &c.resourceProperties.UserPoolID,
			ProviderDetails:  c.resourceProperties.ProviderDetails,
			AttributeMapping: c.resourceProperties.AttributeMapping,
		}).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to update Identity Provider. Error %s", err.Error())
	}

	return map[string]string{
		"ProviderName": *resp.IdentityProvider.ProviderName,
		"ProviderType": string(resp.IdentityProvider.ProviderType),
		"UserPoolId":   *resp.IdentityProvider.UserPoolId,
	}, nil
}
