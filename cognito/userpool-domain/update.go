package main

import (
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// updateDomain updates the domain specified in req.
// Returns map[string]string and error.
func (c *config) updateDomain(req *events.Request) (map[string]string, error) {
	// Due to the API for UserPool Domain being so buggy we need to delete and create.
	if err := c.deleteDomain(req); err != nil {
		return nil, err
	}
	return c.createDomain(req)
}
