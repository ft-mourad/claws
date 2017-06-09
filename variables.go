package main

import "flag"

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
	region        string
	instances     []SimpleInstance
	iids          []string
	regionTag     = flag.String("r", "eu-central-1", "search tag, region")
	Tag           = flag.String("t", "Name:*", "search tag, key")
	keyTag        string
	valueTag      string
	instanceId           = flag.String("i", "", "specific instance id")
	command              = flag.String("c", "", "start or stop command on all instances retrieved by the search")
	instance_mode string = "false"
	command_param []string
)
