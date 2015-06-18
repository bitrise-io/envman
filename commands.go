package main

import "github.com/codegangsta/cli"

var (
	commands = []cli.Command{
		{
			Name:      "init",
			ShortName: "i",
			Usage:     "Create an empty ENVSTORE into ENVMAN_WORK_DIR (i.e. create ENVMAN_WORK_PATH)",
			Action:    initCmd,
		},
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Add new or update an exist environment variable",
			Flags: []cli.Flag{
				flKey,
				flValue,
				flValueFile,
			},
			Action: addCmd,
		},
		{
			Name:      "clear",
			ShortName: "c",
			Usage:     "Clears the envman provided enviroment variables",
			Action:    clearCmd,
		},
		{
			Name:      "print",
			ShortName: "p",
			Usage:     "Prints out the environment variables in ENVMAN_WORK_PATH",
			Action:    printCmd,
		},
		{
			Name:            "run",
			ShortName:       "r",
			Usage:           "Runs the specified command with environment variables in ENVSTORE",
			SkipFlagParsing: true,
			Action:          runCmd,
		},
	}
)
