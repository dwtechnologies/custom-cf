package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// deleteTags will delete all tags for ResourceArn with
// settings specified by req.
// Returns a map of properties and error.
func (c *config) deleteTags(req *events.Request) error {
	tagKeys := []string{}
	for _, tag := range c.oldResourceProperties.Tags {
		tagKeys = append(tagKeys, *tag.Key)
	}

	_, err := c.svc.UntagResourceRequest(
		&ecs.UntagResourceInput{
			ResourceArn: &c.resourceProperties.ResourceARN,
			TagKeys:     tagKeys,
		}).Send()
	if err != nil {
		return fmt.Errorf("Failed to tag resource. Error %s", err.Error())
	}

	return nil
}
