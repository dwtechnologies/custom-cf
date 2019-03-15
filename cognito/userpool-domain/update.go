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
	}

	// The update part of the SDK doesn't work...
	if err := c.deleteDomain(req); err != nil {
		return nil, err
	}

	return c.createDomain(req)
}
