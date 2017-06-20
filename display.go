package main

import (
	"encoding/json"
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
	fmt.Printf("%-40s\t %-60s\t %-60s\t %s\n", SI.Id, SI.Name, SI.Owner, SI.State)
}

func simpleOutput() {
	for _, instance := range instances {
		if instance.State == "running" {
			color.Set(color.FgGreen)
		} else if instance.State == "stopped" {
			color.Set(color.FgRed)
		} else {
			color.Set(color.FgYellow)
		}
		displaySimpleInstances(instance)
	}
	defer color.Unset()
}

func jsonOutput() {
	jsonString := jsonConvert()
	fmt.Println(jsonString)
}

func displayResults() {

	if *format == "json" {
		jsonOutput()
	} else {
		simpleOutput()
	}
}

func indexResult(resp *ec2.DescribeInstancesOutput) {
	var instance SimpleInstance
	//looping through the results (list of instances)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			//create a SimpleInstance object per returned existing EC2 instance
			instance.State = *inst.State.Name
			instance.Id = *inst.InstanceId
			//parse the list of tags to find the Owner and Name tags
			for i := 0; i < len(inst.Tags); i++ {
				if *inst.Tags[i].Key == "Name" {
					instance.Name = *inst.Tags[i].Value
				}
				if *inst.Tags[i].Key == "Owner" {
					instance.Owner = *inst.Tags[i].Value
				}
			}
			//add the Instance ID to an array of IDs
			iids = append(iids, instance.Id)
			//add the SimpleInstance Object to an array of SimpleInstances
			instances = append(instances, instance)
		}
	}
}

func jsonConvert() string {
	jsonString, _ := json.Marshal(instances)
	//fmt.Println(string(jsonString))
	return string(jsonString)
}
