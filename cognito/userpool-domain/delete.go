package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteDomain will delete the domain specified in req.
// If old is true we delete from oldResourceProperties.
// Returns error.
func (c *config) deleteDomain(req *events.Request, old bool) error {
	props := c.resourceProperties
	if old {
		props = c.oldResourceProperties
	}

	switch {
	case props.Domain == "":
		return fmt.Errorf("No Domain specified")

	case props.UserPoolID == "":
		return fmt.Errorf("No UserPoolId specified")
	}

	_, err := c.svc.DeleteUserPoolDomainRequest(
		&cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     &props.Domain,
			UserPoolId: &props.UserPoolID,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to deleteDomain. Error %s", err.Error())
	}

	return nil
}
