package main

import (
	"log"

	"github.com/codegangsta/cli"
)

func runCommand(c *cli.Context) {
	environments, err := loadEnvMap()
	if err != nil {
		log.Fatalln("Failed to export environment variable list, err:", err)
	}

	doCmdEnvs := environments
	doCommand := c.Args()[0]
	doArgs := c.Args()[1:]

	cmdToSend := commandModel{
		Command:      doCommand,
		Environments: doCmdEnvs,
		Argumentums:  doArgs,
	}

	executeCmd(cmdToSend)
}
