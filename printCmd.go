package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func printEnvs() error {
	environments, err := loadEnvMap()
	if err != nil {
		return err
	}

	log.Info("EnvStore:", environments)
	return nil
}

func printCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	err := printEnvs()
	if err != nil {
		log.Fatal("Failed to print:", err)
	}
}
