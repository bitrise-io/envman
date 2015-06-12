package main

import "github.com/codegangsta/cli"

var (
	flPath = cli.StringFlag{
		Name:  "path, p",
		Value: "",
		Usage: "envman's working path, this is file path, with format {SOME_DIR/envstore.yml}",
	}
	flags = []cli.Flag{
		flPath,
	}

	flKey = cli.StringFlag{
		Name:  "key, k",
		Value: "",
		Usage: "key of the environment variable",
	}
	flValue = cli.StringFlag{
		Name:  "value, v",
		Value: "",
		Usage: "value of the environment variable",
	}
)
