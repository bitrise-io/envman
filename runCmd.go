package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func runCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	if len(c.Args()) > 0 {
		doCmdEnvs, err := loadEnvMap()
		if err != nil {
			log.Fatal("Failed to load EnvStore:", err)
		}

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

		log.Info("Executing command:", cmdToExecute)

		if err := executeCmd(cmdToExecute); err != nil {
			log.Fatal("Failed to execute command:", err)
		}

		log.Info("Command executed")
	} else {
		log.Fatal("No command specified")
	}
}
