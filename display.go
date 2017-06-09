package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
)

func addInstanceIDFilter(iid string, param *ec2.DescribeInstancesInput) *ec2.DescribeInstancesInput {
	param = &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(iid),
				},
			},
		},
	}
	return param
}

func displaySimpleInstances(SI SimpleInstance) {
	fmt.Printf("%-40s\t %-60s\t %-60s\t %s\n", SI.id, SI.name, SI.owner, SI.state)
}

func displayResults() {
	for _, instance := range instances {
		if instance.state == "running" {
			color.Set(color.FgGreen)
		} else if instance.state == "stopped" {
			color.Set(color.FgRed)
		} else {
			color.Set(color.FgYellow)
		}
		displaySimpleInstances(instance)
	}
}

func formatResult(resp *ec2.DescribeInstancesOutput) {
	//var tags = make(map[string]string)
	var instance SimpleInstance
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			instance.state = *inst.State.Name
			instance.id = *inst.InstanceId
			for i := 0; i < len(inst.Tags); i++ {
				if *inst.Tags[i].Key == "Name" {
					instance.name = *inst.Tags[i].Value
				}
				if *inst.Tags[i].Key == "Owner" {
					instance.owner = *inst.Tags[i].Value
				}
			}
			iids = append(iids, instance.id)
			instances = append(instances, instance)
		}
	}
	displayResults()
	defer color.Unset()
}
