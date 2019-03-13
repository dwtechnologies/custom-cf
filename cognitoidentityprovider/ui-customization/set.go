package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// setUICustomization will set UI details for clientID on the user pool with
// settings specified by req. If the user pool already exists it will be updated with
// the settings in req.
// Returns error.
func (c *config) setUICustomization(req *events.Request) (map[string]string, error) {
	resp, err := c.svc.SetUICustomizationRequest(
		&cognitoidentityprovider.SetUICustomizationInput{
			CSS:        &c.resourceProperties.CSS,
			ClientId:   &c.resourceProperties.ClientID,
			ImageFile:  c.resourceProperties.ImageFile,
			UserPoolId: &c.resourceProperties.UserPoolID,
		}).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to set UI Customization. Error %s", err.Error())
	}

	return map[string]string{
		"CSSVersion": *resp.UICustomization.CSSVersion,
		"ClientId":   *resp.UICustomization.ClientId,
		"UserPoolId": *resp.UICustomization.UserPoolId,
	}, nil
}
