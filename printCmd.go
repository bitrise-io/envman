package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func printCmd(c *cli.Context) {
	environments, err := loadEnvMap()
	if err != nil {
		log.Fatalln("Failed to print:", err)
	}

	log.Info("EnvStore:", environments)
}
