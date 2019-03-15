package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/dwtechnologies/custom-cf/lib/events"
)

// createTags will set tags for RoleName with
// settings specified by req.
// Returns a map of properties and error.
func (c *config) createTags(req *events.Request) (map[string]string, error) {
	// append CF stack-id
	c.resourceProperties.Tags = append(c.resourceProperties.Tags, iam.Tag{
		Key:   aws.String("cloudformation:stack-id"),
		Value: aws.String(req.StackID),
	})

	// append CF stack-name
	stackName := strings.Split(req.StackID, "/")[1]
	c.resourceProperties.Tags = append(c.resourceProperties.Tags, iam.Tag{
		Key:   aws.String("cloudformation:stack-name"),
		Value: aws.String(stackName),
	})

	_, err := c.svc.TagRoleRequest(
		&iam.TagRoleInput{
			RoleName: &c.resourceProperties.RoleName,
			Tags:     c.resourceProperties.Tags,
		}).Send()
	if err != nil {
		return nil, fmt.Errorf("Failed to tag role. Error %s", err.Error())
	}

	return map[string]string{}, nil
}
