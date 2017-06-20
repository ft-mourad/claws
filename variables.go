package main

import "flag"

type SimpleInstance struct {
	Id    string
	Name  string
	State string
	Owner string
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
	web           = flag.Bool("w", false, "-w to start web interface")
	keyTag        string
	valueTag      string
	instanceId    = flag.String("i", "", "specific instance id")
	command       = flag.String("c", "", "start or stop command on all instances retrieved by the search")
	format        = flag.String("f", "text", "change output format {text;json}")
	command_param []string
)
