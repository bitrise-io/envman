package main

import "github.com/codegangsta/cli"

var (
	flPath = cli.StringFlag{
		Name:  "path",
		Value: "",
		Usage: "path of the environment variables",
	}
	flags = []cli.Flag{
		flPath,
	}

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
