package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var runCmdLog *log.Entry = log.WithFields(log.Fields{"f": "runCmd.go"})

func runCmd(c *cli.Context) {
	environments, err := loadEnvMap()
	if err != nil {
		runCmdLog.Fatal("Failed to export environment variable list, err:", err)
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
			runCmdLog.Fatal("Failed to execute command, err:", err)
		}
	} else {
		runCmdLog.Fatal("No command specified")
	}
}
