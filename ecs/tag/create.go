package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// createTags will set tags for ResourceArn with
// settings specified by req.
// Returns a map of properties and error.
func (c *config) createTags(req *events.Request) (map[string]string, error) {
	_, err := c.svc.TagResourceRequest(
		&ecs.TagResourceInput{
			ResourceArn: &c.resourceProperties.ResourceARN,
			Tags:        c.resourceProperties.Tags,
		}).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to tag resource. Error %s", err.Error())
	}

	return map[string]string{}, nil
}
