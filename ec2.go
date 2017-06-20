package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ft-mourad/libSimpleEC2"
)

func parseNewTags(svc *ec2.EC2, iid string) {
	for _, param := range command_param {
		parsed_arg := strings.Split(param, ":")
		keyTag, valueTag = parsed_arg[0], parsed_arg[1]
		SEC2.TagInstance(svc, iid, keyTag, valueTag)
	}
}
