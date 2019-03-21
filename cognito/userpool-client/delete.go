package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteClient will delete the UserPool Client specified with clientID.
// Returns error.
func (c *config) deleteClient(req *events.Request, id string) error {
	// If resource creation fails. We will get an empty delete event.
	// In these cases just return nil.
	switch {
	case id == "":
		return nil

	case c.resourceProperties.UserPoolID == "":
		return nil
	}

	_, err := c.svc.DeleteUserPoolClientRequest(
		&cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   &id,
			UserPoolId: &c.resourceProperties.UserPoolID,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to delete Client. Error %s", err.Error())
	}

	return nil
}
