package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func runCmd(c *cli.Context) {
	doCmdEnvs, err := loadEnvMap()
	if err != nil {
		log.Fatal("Failed to load EnvStore:", err)
	}

	if len(c.Args()) > 0 {
		doCommand := c.Args()[0]

		doArgs := []string{}
		if len(c.Args()) > 1 {
			doArgs = c.Args()[1:]
		}

		cmdToExecute := commandModel{
			Command:      doCommand,
			Environments: doCmdEnvs,
			Argumentums:  doArgs,
		}

		log.Debugln("Executing comand:", cmdToExecute)

		if err := executeCmd(cmdToExecute); err != nil {
			log.Fatal("Failed to execute comand:", err)
		}

		log.Info("Comand executed")
	} else {
		log.Fatal("No comand specified")
	}
}
