package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// updateTags will set tags for ResourceArn with
// settings specified by req.
// Returns a map of properties and error.
func (c *config) updateTags(req *events.Request) (map[string]string, error) {
	// get current tags
	curTagKeys := []string{}
	for _, tag := range c.oldResourceProperties.Tags {
		curTagKeys = append(curTagKeys, *tag.Key)
	}

	// get new tags
	newTagKeys := []string{}
	for _, tag := range c.resourceProperties.Tags {
		newTagKeys = append(newTagKeys, *tag.Key)
	}

	// remove tags not current
	unTags := slicesDiff(curTagKeys, newTagKeys)
	if len(unTags) > 0 {
		_, err := c.svc.UntagResourceRequest(
			&ecs.UntagResourceInput{
				ResourceArn: &c.resourceProperties.ResourceARN,
				TagKeys:     unTags,
			}).Send()
		if err != nil {
			return nil, fmt.Errorf("Failed to tag resource. Error %s", err.Error())
		}
	}

	// add c.resourceProperties.Tags tags
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

func slicesDiff(s, k []string) []string {
	diff := []string{}
	m := map[string]int{}

	// first slice map
	for _, val := range s {
		m[val] = 1
	}

	// second slice map
	for _, val := range k {
		m[val] = m[val] + 1
	}

	for key, val := range m {
		if val == 1 {
			diff = append(diff, key)
		}
	}

	return diff
}