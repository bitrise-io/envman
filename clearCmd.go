package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

func clearEnvs() error {
	isExists, err := pathutil.IsPathExists(currentEnvStoreFilePath)
	if err != nil {
		return err
	}
	if !isExists {
		errMsg := "EnvStore not found in path:" + currentEnvStoreFilePath
		return errors.New(errMsg)
	}

	err = writeEnvMapToFile(currentEnvStoreFilePath, []envModel{})
	if err != nil {
		return err
	}

	return nil
}

func clearCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	err := clearEnvs()
	if err != nil {
		log.Fatal("Failed to clear EnvStore:", err)
	}

	log.Info("EnvStore cleared")
}
