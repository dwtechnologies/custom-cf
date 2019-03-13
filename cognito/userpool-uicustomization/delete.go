package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

const defaultCSS = ".logo-customizable {max-width: 60%;max-height: 30%;}.banner-customizable {padding: 25px 0px 25px 0px;background-color: lightgray;}.label-customizable {font-weight: 400;}.textDescription-customizable {padding-top: 10px;padding-bottom: 10px;display: block;font-size: 16px;}.idpDescription-customizable {padding-top: 10px;padding-bottom: 10px;display: block;font-size: 16px;}.legalText-customizable {color: #747474;font-size: 11px;}.submitButton-customizable {font-size: 14px;font-weight: bold;margin: 20px 0px 10px 0px;height: 40px;width: 100%;color: #fff;background-color: #337ab7;}.submitButton-customizable:hover {color: #fff;background-color: #286090;}.errorMessage-customizable {padding: 5px;font-size: 14px;width: 100%;background: #F5F5F5;border: 2px solid #D64958;color: #D64958;}.inputField-customizable {width: 100%;height: 34px;color: #555;background-color: #fff;border: 1px solid #ccc;}.inputField-customizable:focus {border-color: #66afe9;outline: 0;}.idpButton-customizable {height: 40px;width: 100%;text-align: center;margin-bottom: 15px;color: #fff;background-color: #5bc0de;border-color: #46b8da;}.idpButton-customizable:hover {color: #fff;background-color: #31b0d5;}.socialButton-customizable {height: 40px;text-align: left;width: 100%;margin-bottom: 15px;}.redirect-customizable {text-align: center;}.passwordCheck-notValid-customizable {color: #DF3312;}.passwordCheck-valid-customizable {color: #19BF00;}.background-customizable {background-color: #fff;}"

// deleteUICustomization will resetset UI details for clientID on the user pool
// Returns error.
func (c *config) deleteUICustomization(req *events.Request) error {
	_, err := c.svc.SetUICustomizationRequest(
		&cognitoidentityprovider.SetUICustomizationInput{
			CSS:      aws.String(defaultCSS),
			ClientId: &c.resourceProperties.ClientID,
			// ImageFile:  c.resourceProperties.ImageFile,
			UserPoolId: &c.resourceProperties.UserPoolID,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to delete UI Customization. Error %s", err.Error())
	}

	return nil
}
