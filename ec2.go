package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func EC2_init(region string) *ec2.EC2 {
	//create session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	//initialize session
	svc := ec2.New(sess, &aws.Config{Region: aws.String(region)})
	return svc
}

func addTagFilter(keyTag string, valueTag string, param *ec2.DescribeInstancesInput) *ec2.DescribeInstancesInput {
	param = &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + keyTag),
				Values: []*string{
					aws.String(valueTag),
				},
			},
		},
	}
	return param
}

func searchInstances(svc *ec2.EC2) {
	var inputFilter *ec2.DescribeInstancesInput
	//define inputs (filters)
	//The type DescribeInstancesInput obviously contains the inputs (here filter)
	//for the DescribeInstances below
	switch instance_mode {
	case "false":
		inputFilter = addTagFilter(keyTag, valueTag, inputFilter)
	case "true":
		inputFilter = addInstanceIDFilter(*instanceId, inputFilter)
	}
	//call
	resp, err := svc.DescribeInstances(inputFilter)
	if err != nil {
		panic(err)
	}
	formatResult(resp)
}

func createNewTags(svc *ec2.EC2, iid string, keyTag string, valueTag string) {
	fmt.Println("create tag -- ", keyTag, " : ", valueTag, "\n")
	svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(iid)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(keyTag),
				Value: aws.String(valueTag),
			},
		},
	})
}

func parseNewTags(svc *ec2.EC2, iid string) {
	for _, param := range command_param {
		parsed_arg := strings.Split(param, ":")
		keyTag, valueTag = parsed_arg[0], parsed_arg[1]
		createNewTags(svc, iid, keyTag, valueTag)
	}
}
