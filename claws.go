package main

// This is just an experimentationm, just a test to learn go and familiarize
// myself with the aws sdk
//
//What I am trying to achieve here : a simple tag-based AWS ECS admin tool
//Just focused on what I need
//This will not do anythging complex, simply:
//  search and list instance matching a specified tag
//  enable easy start/stop (no terminate) operation on id or tags
//  maybe link to the billing services

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
)

var (
	regions = map[string]string{
		"ireland":   "eu-west-1",
		"frankfurt": "eu-central-1",
	}

	iids          []string
	regionTag     = flag.String("r", "eu-central-1", "search tag, region")
	Tag           = flag.String("t", "Name:*", "search tag, key")
	keyTag        string
	valueTag      string
	instanceId           = flag.String("i", "", "specific instance id")
	instance_mode string = "false"
)

// CLI arguments

func parseInput() string {

	//process the arguments
	var region string

	flag.Parse()
	if tmp, exist := regions[*regionTag]; exist {
		region = tmp
	} else {
		region = *regionTag
	}

	parsed_arg := strings.Split(*Tag, ":")
	keyTag, valueTag = parsed_arg[0], parsed_arg[1]
	fmt.Println(keyTag, " : ", valueTag)

	if *instanceId != "" && keyTag != "" {
		fmt.Println(" -i and -k/v are mutualy exclusive. Please use one or the other\n")
		instance_mode = "true"
	}

	fmt.Println("instance : ", *instanceId, "\n")
	fmt.Println("search : ", keyTag, " : ", valueTag, " in ", *regionTag, "\n")

	return region
}

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

// func stopInstances(iids []string){
//
//   var param *ec2.StopInstancesInput
//   param = param.SetInstanceIDs(iids []string)
// }

func formatResult(resp *ec2.DescribeInstancesOutput) {
	var tags = make(map[string]string)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			//var pos int = 0
			//Here, i couldn't figure out how to find the Name tag on the instance
			//I just assume if it is not at the 1st position it is at the second
			//very ugly
			//if *inst.Tags[0].Key == "Name" {
			// } else {
			// 	pos = 1
			// }

			//to make things nicer i use fatih/color" to colorize the output
			//depending on the status
			if *inst.State.Name == "running" {
				color.Set(color.FgGreen)
			} else if *inst.State.Name == "stopped" {
				color.Set(color.FgRed)
			}

			for i := 0; i <= len(inst.Tags)-1; i++ {
				tags[*inst.Tags[i].Key] = *inst.Tags[i].Value
			}
			fmt.Println(tags["Name"], "\t\t")
			//fmt.Println(*inst.Tags[0].Key, len(inst.Tags))
			//iids = append(iids, *inst.InstanceId)
			//fmt.Printf("%-40s\t %-60s\t %s\n", *inst.InstanceId, *inst.Tags[pos].Value, *inst.State.Name)

		}
	}
	//fmt.Println(iids)
	defer color.Unset() // Use it in your function

}

func searchInstances(svc) {
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
	resp, err := svc.DescribeInstances(inputFilter)
	if err != nil {
		panic(err)
	}

	formatResult(resp)
}

//need an index func to put into slices/maps the search result.
//then we'll be able to run command against the Instances
func main() {

	region := parseInput()
	svc := EC2_init(region)

	searchInstances(svc)
	//run the 'list Instances command'
	//resp is of type DescribeInstancesOutput

}
