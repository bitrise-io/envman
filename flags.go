package main

import "github.com/codegangsta/cli"

var (
	flKey = cli.StringFlag{
		Name:  "key",
		Value: "",
		Usage: "key of the environment variable",
	}
	flValue = cli.StringFlag{
		Name:  "value",
		Value: "",
		Usage: "value of the environment variable",
	}
)
