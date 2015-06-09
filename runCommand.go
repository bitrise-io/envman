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

	if len(c.Args()) > 0 {
		doCmdEnvs := environments
		doCommand := c.Args()[0]

		doArgs := []string{}
		if len(c.Args()) > 1 {
			doArgs = c.Args()[1:]
		}

		cmdToSend := commandModel{
			Command:      doCommand,
			Environments: doCmdEnvs,
			Argumentums:  doArgs,
		}

		if err := executeCmd(cmdToSend); err != nil {
			log.Fatalln("Failed to execute command, err:", err)
		}
	} else {
		log.Fatalln("No command specified")
	}
}
