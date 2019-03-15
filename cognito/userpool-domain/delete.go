package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteDomain will delete the domain specified in req.
// Returns error.
func (c *config) deleteDomain(req *events.Request) error {
	switch {
	case c.resourceProperties.Domain == "":
		return fmt.Errorf("No Domain specified")

	case c.resourceProperties.UserPoolID == "":
		return fmt.Errorf("No UserPoolId specified")
	}

	_, err := c.svc.DeleteUserPoolDomainRequest(
		&cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     &c.resourceProperties.Domain,
			UserPoolId: &c.resourceProperties.UserPoolID,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to deleteDomain. Error %s", err.Error())
	}

	return nil
}
