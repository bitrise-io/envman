package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var clearCmdLog *log.Entry = log.WithFields(log.Fields{"f": "clearCmd.go"})

func clearCmd(c *cli.Context) {
	err := writeEnvMapToFile(currentEnvStoreFilePath, envMap{})
	if err != nil {
		clearCmdLog.Error("Failed to clear envlist, err:%s", err)
	}
}
