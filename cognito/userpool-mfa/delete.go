package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteMFA will delete the MFA settings and set them to default values.
// Returns error.
func (c *config) deleteMFA(req *events.Request) error {
	switch {
	case c.Domain == "":
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
		return fmt.Errorf("Failed to delete Domain. Error %s", err.Error())
	}

	return nil
}
