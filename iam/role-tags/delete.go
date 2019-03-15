package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteTags will delete all tags for RoleName with
// settings specified by req.
// Returns a map of properties and error.
func (c *config) deleteTags(req *events.Request) error {
	// get current tags
	curTagKeys := []string{}
	for _, tag := range c.resourceProperties.Tags {
		curTagKeys = append(curTagKeys, *tag.Key)
	}

	// append CF stack-id, stack-name
	curTagKeys = append(curTagKeys, "cloudformation:stack-id")
	curTagKeys = append(curTagKeys, "cloudformation:stack-name")

	_, err := c.svc.UntagRoleRequest(
		&iam.UntagRoleInput{
			RoleName: &c.resourceProperties.RoleName,
			TagKeys:  curTagKeys,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to tag role. Error %s", err.Error())
	}

	return nil
}
