package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func runWithEnvs(cmd commandModel) error {
	if err := executeCmd(cmd); err != nil {
		return err
	}

	return nil
}

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

		log.Info("Executing comand:", cmdToExecute)

		err = runWithEnvs(cmdToExecute)
		if err != nil {
			log.Fatal("Failed to execute comand:", err)
		}

		log.Info("Comand executed")
	} else {
		log.Fatal("No comand specified")
	}
}
