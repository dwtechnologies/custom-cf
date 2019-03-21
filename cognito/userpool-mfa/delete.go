package main

import "github.com/dwtechnologies/custom-cf/lib/events"

// deleteMFA will delete the MFA configuration for the UserPool.
// Returns error.
func (c *config) deleteMFA(req *events.Request) error {
	// If UserPoolID is empty we must assume that an delete was sent
	// on a failed resource creation. So just return nil.
	switch {
	case c.resourceProperties.UserPoolID == "":
		return nil
	}

	return c.setMFA(req, true)
}
