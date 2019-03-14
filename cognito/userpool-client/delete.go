package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteClient will delete the UserPool Client specified with clientID.
// Returns error.
func (c *config) deleteClient(req *events.Request, id string) error {
	switch {
	case id == "":
		return fmt.Errorf("Client ID can't be empty")

	case c.resourceProperties.UserPoolID == "":
		return fmt.Errorf("UserPool ID can't be empty")
	}

	_, err := c.svc.DeleteUserPoolClientRequest(
		&cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   &c.resourceProperties.id,
			UserPoolId: &c.resourceProperties.UserPoolID,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to delete Client. Error %s", err.Error())
	}

	return nil
}
