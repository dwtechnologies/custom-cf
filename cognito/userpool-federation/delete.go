package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteIdentityProvider will delete the Identity Provider specified in req.
// Returns error.
func (c *config) deleteIdentityProvider(req *events.Request) error {
	// We will get an empty delete event if resource creation fails.
	// So we need to return nil on these events.
	switch {
	case c.resourceProperties.UserPoolID == "":
		return nil

	case c.resourceProperties.ProviderName == "":
		return nil
	}

	_, err := c.svc.DeleteIdentityProviderRequest(
		&cognitoidentityprovider.DeleteIdentityProviderInput{
			ProviderName: &c.resourceProperties.ProviderName,
			UserPoolId:   &c.resourceProperties.UserPoolID,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to delete Identity Provider. Error %s", err.Error())
	}

	return nil
}
