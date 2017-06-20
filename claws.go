//What I am trying to achieve here : a simple tag-based AWS ECS admin tool
//Just focused on what I need

//This tool will be considered feature-full when it can :
// - List exiting instances on any region, their name, Owner tag and their status
// - Search instances by ID or Tags
// - Start and Stop command AGAINST ALL INSTANCES MATCHING THE SEARCH FILTER
// - Add tags or retags command AGAINST ALL INSTANCES MATCHING THE SEARCH FILTER

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ft-mourad/libSimpleEC2"
)

//ran once
func main() {
	parseInput()
	svc := SEC2.EC2_init(region)
	_, iids := SEC2.ListInstances(svc)
	commandInstances(svc, iids)
}

// CLI arguments
// CLAWS expect any of the following option as arguments
// Region: -r REGION      : could be eu-west-1 or ireland format
// Tags:   -t KEY:VALUE   : Tag-based search, wildcard can be used
// Command: -c start / -c stop /-c tag Owner:Batman : add or edit existing tag
func parseInput() {
	//process the arguments
	flag.Parse()
	command_param = flag.Args()
	if tmp, exist := regions[*regionTag]; exist {
		region = tmp
	} else {
		region = *regionTag
	}
	parsed_arg := strings.Split(*Tag, ":")
	keyTag, valueTag = parsed_arg[0], parsed_arg[1]

}

//As each command ran has different config/syntax, below is the parsing.
//For each instance found in the tag-based search
//Build AWS object (ionput/output)
//Create resources
func commandInstances(svc *ec2.EC2, iids []string) {
	//loop through all instances
	if *command != "" {
		for _, iid := range iids {
			// // start commamd
			// if *command == "start" {
			// 	fmt.Println("STARTATATATATAT")
			// 	SEC2.StartInstance(svc, iid)
			// 	//stop command
			// } else if *command == "stop" {
			// 	fmt.Println("STOPOPOPOPOPOPO")
			// 	SEC2.StopInstance(svc, iid)
			// } else if *command == "tag" {
			// 	fmt.Println(iid)
			// }

			switch {
			case *command == "start":
				fmt.Println("Starting instance : ", iid)
				SEC2.StartInstance(svc, iid)
			case *command == "stop":
				fmt.Println("Stoping instance : ", iid)
				SEC2.StopInstance(svc, iid)
			case *command == "tag":
				fmt.Println("Editing tag : ", iid)
				parseNewTags(svc, iid)
			default:
				fmt.Println("unsupported command")
			}
		}
	}
}
