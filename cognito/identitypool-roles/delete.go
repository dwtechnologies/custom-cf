package main

import "github.com/dwtechnologies/custom-cf/lib/events"

// deleteRoles will delete the roles configuration for the IdentityPool and set empty roles.
// Returns error.
func (c *config) deleteRoles(req *events.Request) error {
	// If IdentityPoolID is empty we must assume that an delete was sent
	// on a failed resource creation. So just return nil.
	switch {
	case c.resourceProperties.IdentityPoolID == "":
		return nil
	}

	return c.setRoles(req, true)
}
