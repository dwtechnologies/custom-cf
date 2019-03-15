package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// updateDomain updates the domain specified in req.
// Returns map[string]string and error.
func (c *config) updateDomain(req *events.Request) (map[string]string, error) {
	switch {
	case c.resourceProperties.Domain == "":
		return nil, fmt.Errorf("No Domain specified")

	case c.resourceProperties.UserPoolID == "":
		return nil, fmt.Errorf("No UserPoolId specified")
	}

	input := &cognitoidentityprovider.UpdateUserPoolDomainInput{
		Domain:     &c.resourceProperties.Domain,
		UserPoolId: &c.resourceProperties.UserPoolID,
	}

	// Only set CustomDomainConfig if it's not nil.
	if c.resourceProperties.CustomDomainConfig != nil {
		input.CustomDomainConfig = &cognitoidentityprovider.CustomDomainConfigType{
			CertificateArn: &c.resourceProperties.CustomDomainConfig.CertificateArn,
		}
	}

	resp, err := c.svc.UpdateUserPoolDomainRequest(input).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to update Domain. Error %s", err.Error())
	}

	data := map[string]string{}
	if resp.CloudFrontDomain != nil {
		data["Domain"] = *resp.CloudFrontDomain
	}

	return data, nil
}
