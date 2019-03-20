package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// setMFA will set the MFA settings specified by req.
// Returns map[string]string and error.
func (c *config) setMFA(req *events.Request) (map[string]string, error) {
	switch {
	case c.resourceProperties.MfaConfiguration != "OFF" || c.resourceProperties.MfaConfiguration != "ON" || c.resourceProperties.MfaConfiguration != "OPTIONAL":
		return nil, fmt.Errorf("No MfaConfiguration needs to be either OFF ON or OPTIONAL")

	case c.resourceProperties.UserPoolID == "":
		return nil, fmt.Errorf("No UserPoolId specified")
	}

	input := &cognitoidentityprovider.SetUserPoolMfaConfigInput{
		MfaConfiguration: cognitoidentityprovider.UserPoolMfaType(c.resourceProperties.MfaConfiguration),
		UserPoolId:       &c.resourceProperties.UserPoolID,
	}

	// Only set SmsMfaConfiguration if it's not nil.
	if c.resourceProperties.SmsMfaConfiguration != nil {
		input.SmsMfaConfiguration = &cognitoidentityprovider.SmsMfaConfigType{
			SmsAuthenticationMessage: &c.resourceProperties.SmsMfaConfiguration.SmsAuthenticationMessage,
		}
		if c.resourceProperties.SmsMfaConfiguration.SmsConfiguration != nil {
			input.SmsMfaConfiguration.SmsConfiguration = &cognitoidentityprovider.SmsConfigurationType{
				SnsCallerArn: &c.resourceProperties.SmsMfaConfiguration.SmsConfiguration.SnsCallerArn,
				ExternalId:   &c.resourceProperties.SmsMfaConfiguration.SmsConfiguration.ExternalID,
			}
		}
	}

	// Only set SoftwareMfaConfiguration if it's not nil.
	if c.resourceProperties.SoftwareTokenMfaConfiguration != nil {
		input.SoftwareTokenMfaConfiguration = &cognitoidentityprovider.SoftwareTokenMfaConfigType{
			Enabled: &c.resourceProperties.SoftwareTokenMfaConfiguration.Enabled,
		}
	}

	_, err := c.svc.SetUserPoolMfaConfigRequest(input).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to set MFA. Error %s", err.Error())
	}

	return nil, nil
}
