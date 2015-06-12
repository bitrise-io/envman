package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

var clearCmdLog *log.Entry = log.WithFields(log.Fields{"f": "clearCmd.go"})

func clearCmd(c *cli.Context) {
	isExists, err := pathutil.IsPathExists(currentEnvStoreFilePath)
	if err != nil {
		clearCmdLog.Error("Failed to clear envlist, err:%s", err)
		return
	}
	if !isExists {
		clearCmdLog.Info("No EnvStore found at path: ", currentEnvStoreFilePath)
		return
	}

	err = writeEnvMapToFile(currentEnvStoreFilePath, envMap{})
	if err != nil {
		clearCmdLog.Error("Failed to clear envlist, err:%s", err)
		return
	}
	clearCmdLog.Info("Envstore cleared at path: ", currentEnvStoreFilePath)
}
