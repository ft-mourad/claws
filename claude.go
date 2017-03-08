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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
)

// TO DO
//    Add error-catching (this is just a quick coding session)
//    Restructure in functions
//    Find a nicer way to parse Tags[]
//
//    ... and add more stuff

// CLI arguments
var regionTag = flag.String("r", "eu-central-1", "search tag, region")
var keyTag = flag.String("k", "Name", "search tag, key")
var valueTag = flag.String("v", "*", "search tag, value")

func main() {

	//process the arguments
	flag.Parse()
	fmt.Println("search : ", *keyTag, " : ", *valueTag, " in ", *regionTag, "\n")

	//create session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	//initialize session
	svc := ec2.New(sess, &aws.Config{Region: aws.String(*regionTag)})

	//define inputs (filters)
	//The type DescribeInstancesInput obviously contains the inputs (here filter)
	//for the DescribeInstances below
	params := &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(*valueTag),
				},
			},
		},
	}

	//run the 'list Instances command'
	//resp is of type DescribeInstancesOutput
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	//displaying the DescribeInstancesOutput
	//there are arrays inside arrays here, so nested FORs
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			var pos int = 0
			//Here, i couldn't figure out how to find the Name tag on the instance
			//I just assume if it is not at the 1st position it is at the second
			//very ugly
			if *inst.Tags[0].Key == "Name" {
			} else {
				pos = 1
			}

			//to make things nicer i use atih/color" to colorize the output
			//depending on the status
			if *inst.State.Name == "running" {
				color.Set(color.FgGreen)
			} else if *inst.State.Name == "stopped" {
				color.Set(color.FgRed)
			}

			fmt.Printf("%-20s\t %-60s\t %s\n", *inst.InstanceId, *inst.Tags[pos].Value, *inst.State.Name)

		}
	}
	defer color.Unset() // Use it in your function

}
