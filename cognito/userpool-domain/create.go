package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// createDomain will create a new domain on the user pool with
// settings specified by req.
// Returns ap[string]string and error.
func (c *config) createDomain(req *events.Request) (map[string]string, error) {
	switch {
	case c.resourceProperties.Domain == "":
		return nil, fmt.Errorf("No Domain specified")

	case c.resourceProperties.UserPoolID == "":
		return nil, fmt.Errorf("No UserPoolId specified")
	}

	input := &cognitoidentityprovider.CreateUserPoolDomainInput{
		Domain:     &c.resourceProperties.Domain,
		UserPoolId: &c.resourceProperties.UserPoolID,
	}

	// Only set CustomDomainConfig if it's not nil.
	if c.resourceProperties.CustomDomainConfig != nil {
		input.CustomDomainConfig = &cognitoidentityprovider.CustomDomainConfigType{
			CertificateArn: &c.resourceProperties.CustomDomainConfig.CertificateArn,
		}
	}

	resp, err := c.svc.CreateUserPoolDomainRequest(input).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to create Domain. Error %s", err.Error())
	}

	return map[string]string{
		"Domain": *resp.CloudFrontDomain,
	}, nil
}
