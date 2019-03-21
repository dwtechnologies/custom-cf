package main

import (
	"fmt"

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

	case c.oldResourceProperties.Domain == "":
		return nil, fmt.Errorf("No Old Domain specified")

	case c.oldResourceProperties.UserPoolID == "":
		return nil, fmt.Errorf("No Old UserPoolId specified")
	}

	// Due to the API for UserPool Domain being so buggy we need to delete and create.
	if err := c.deleteDomain(req, true); err != nil {
		return nil, err
	}
	return c.createDomain(req)
}
