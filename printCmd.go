package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var printCmdLog *log.Entry = log.WithFields(log.Fields{"f": "printCmd.go"})

func printCmd(c *cli.Context) {
	environments, err := loadEnvMap()
	if err != nil {
		printCmdLog.Fatal("Failed to print environment variable list, err:", err)
	}
	if len(environments) == 0 {
		printCmdLog.Info("Empty environment variable list")
	} else {
		printCmdLog.Info(environments)
	}
}
