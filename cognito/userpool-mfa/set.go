package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// setMFA will set the MFA settings specified by req.
// If defaults is true the default MFA settings will be set.
// error.
func (c *config) setMFA(req *events.Request, defaults bool) error {
	props := c.resourceProperties
	if defaults {
		props = &MFA{MfaConfiguration: "OFF"}
		props.UserPoolID = c.resourceProperties.UserPoolID
	}

	switch {
	case props.MfaConfiguration != "OFF" && props.MfaConfiguration != "ON" && props.MfaConfiguration != "OPTIONAL":
		return fmt.Errorf("No MfaConfiguration needs to be either OFF, ON or OPTIONAL")

	case props.UserPoolID == "":
		return fmt.Errorf("No UserPoolId specified")
	}

	input := &cognitoidentityprovider.SetUserPoolMfaConfigInput{
		MfaConfiguration: cognitoidentityprovider.UserPoolMfaType(props.MfaConfiguration),
		UserPoolId:       &props.UserPoolID,
	}

	// Only set SmsMfaConfiguration if it's not nil.
	if props.SmsMfaConfiguration != nil {
		input.SmsMfaConfiguration = &cognitoidentityprovider.SmsMfaConfigType{
			SmsAuthenticationMessage: &props.SmsMfaConfiguration.SmsAuthenticationMessage,
		}
		if props.SmsMfaConfiguration.SmsConfiguration != nil {
			input.SmsMfaConfiguration.SmsConfiguration = &cognitoidentityprovider.SmsConfigurationType{
				SnsCallerArn: &props.SmsMfaConfiguration.SmsConfiguration.SnsCallerArn,
				ExternalId:   &props.SmsMfaConfiguration.SmsConfiguration.ExternalID,
			}
		}
	}

	// Only set SoftwareMfaConfiguration if it's not nil.
	if props.SoftwareTokenMfaConfiguration != nil {
		val := false
		if strings.ToLower(props.SoftwareTokenMfaConfiguration.Enabled) == "true" {
			val = true
		}

		input.SoftwareTokenMfaConfiguration = &cognitoidentityprovider.SoftwareTokenMfaConfigType{
			Enabled: &val,
		}
	}

	_, err := c.svc.SetUserPoolMfaConfigRequest(input).Send()
	if err != nil {
		return fmt.Errorf("Failed to set MFA. Error %s", err.Error())
	}

	return nil
}
