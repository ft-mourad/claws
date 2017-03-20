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

type SimpleInstance struct {
	id    string
	name  string
	state string
	owner string
}

var (
	regions = map[string]string{
		"ireland":   "eu-west-1",
		"frankfurt": "eu-central-1",
	}
	instances     []SimpleInstance
	iids          []string
	regionTag     = flag.String("r", "eu-central-1", "search tag, region")
	Tag           = flag.String("t", "Name:*", "search tag, key")
	keyTag        string
	valueTag      string
	instanceId           = flag.String("i", "", "specific instance id")
	command              = flag.String("c", "", "start or stop command on all instances retrieved by the search")
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
	fmt.Println("run : ", *command)
	fmt.Println("search : ", *instanceId, "with ", keyTag, " = ", valueTag, " in ", *regionTag, "\n")
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

func commandInstances(svc *ec2.EC2) {

	for _, iid := range iids {
		if *command == "start" {
			fmt.Println("starting")
			input := &ec2.StartInstancesInput{
				InstanceIds: []*string{
					aws.String(iid),
				},
				DryRun: aws.Bool(false),
			}
			_, err := svc.StartInstances(input)
			if err != nil {
				fmt.Println(err)
			}
		} else if *command == "stop" {
			fmt.Println("stopping")
			input := &ec2.StopInstancesInput{
				InstanceIds: []*string{
					aws.String(iid),
				},
				DryRun: aws.Bool(false),
			}
			_, err := svc.StopInstances(input)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
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

func displaySimpleInstances(SI SimpleInstance) {
	fmt.Printf("%-40s\t %-60s\t %-60s\t %s\n", SI.id, SI.name, SI.owner, SI.state)
}

func displayResults() {
	for _, instance := range instances {
		if instance.state == "running" {
			color.Set(color.FgGreen)
		} else if instance.state == "stopped" {
			color.Set(color.FgRed)
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
	//fmt.Println(iids)
	displayResults()
	defer color.Unset() // Use it in your function
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

//need an index func to put into slices/maps the search result.
//then we'll be able to run command against the Instances
func main() {
	region := parseInput()
	svc := EC2_init(region)

	searchInstances(svc)
	commandInstances(svc)
}
