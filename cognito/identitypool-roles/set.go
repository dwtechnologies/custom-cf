package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"

	"github.com/dwtechnologies/custom-cf/lib/events"
)

// setRoles will set the roles specified by req.
// If defaults is true the default roles will be set.
// Returns error.
func (c *config) setRoles(req *events.Request, defaults bool) error {
	props := c.resourceProperties
	if defaults {
		props = &IdentityPoolRoles{Roles: &Roles{Authenticated: "", UnAuthenticated: ""}}
		props.IdentityPoolID = c.resourceProperties.IdentityPoolID
	}

	switch {
	case props.IdentityPoolID == "":
		return fmt.Errorf("No UserPoolId specified")
	}

	// Create the input for the request.
	input := &cognitoidentity.SetIdentityPoolRolesInput{
		Roles: map[string]string{
			"authenticated":   props.Roles.Authenticated,
			"unauthenticated": props.Roles.UnAuthenticated,
		},
		IdentityPoolId: &props.IdentityPoolID,
	}

	// Create the map for RoleMappings if it's not nil.
	if props.RoleMappings != nil {
		input.RoleMappings = map[string]cognitoidentity.RoleMapping{}
	}

	// If RoleMappings is set validate the basic input.
	for _, mapping := range props.RoleMappings {
		switch {
		case mapping.IdentityProvider == "":
			return fmt.Errorf("No IdentityProvider set in RoleMappings")

		case mapping.Type != "Token" && mapping.Type != "Rules":
			return fmt.Errorf("Type is not valid in RoleMappings. Valid values are Token or Rules")

		case mapping.AmbiguousRoleResolution != "AuthenticatedRole" && mapping.AmbiguousRoleResolution != "Deny":
			return fmt.Errorf("AmbiguousRoleResolution is not valid in RoleMappings. Valid values are AuthenticatedRole or Deny")

		case mapping.Type == "Rules" && mapping.Rules == nil:
			return fmt.Errorf("No Rules set in RoleMappings and Type is Rules")
		}

		input.RoleMappings[mapping.IdentityProvider] = cognitoidentity.RoleMapping{
			Type:                    cognitoidentity.RoleMappingType(mapping.Type),
			AmbiguousRoleResolution: cognitoidentity.AmbiguousRoleResolutionType(mapping.AmbiguousRoleResolution),
		}

		// Set the Rules Config struct and map if it's not nil.
		if mapping.Rules != nil {
			r := input.RoleMappings[mapping.IdentityProvider]
			r.RulesConfiguration = &cognitoidentity.RulesConfigurationType{Rules: []cognitoidentity.MappingRule{}}

			// Validate rules.
			for _, rule := range mapping.Rules {
				switch {
				case rule.Claim == "":
					return fmt.Errorf("No Claim set in Rules")

				case rule.MatchType != "Equals" && rule.MatchType != "Contains" && rule.MatchType != "StarsWith" && rule.MatchType != "NotEqual":
					return fmt.Errorf("MatchType is not valid in Rules. Valid values are Equals, Contains, StarsWith or NotEqual")

				case rule.Value == "":
					return fmt.Errorf("No Value set in Rules")

				case rule.RoleArn == "":
					return fmt.Errorf("No RoleArn set in Rules")
				}

				r.RulesConfiguration.Rules = append(r.RulesConfiguration.Rules, cognitoidentity.MappingRule{
					Claim:     &rule.Claim,
					MatchType: cognitoidentity.MappingRuleMatchType(rule.MatchType),
					Value:     &rule.Value,
					RoleARN:   &rule.RoleArn,
				})
			}
		}
	}

	// Send the request.
	_, err := c.svc.SetIdentityPoolRolesRequest(input).Send()
	if err != nil {
		return fmt.Errorf("Failed to set Identity Pool Roles. Error %s", err.Error())
	}

	return nil
}
