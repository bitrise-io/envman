package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

func clearEnvs() error {

	if isExists, err := pathutil.IsPathExists(currentEnvStoreFilePath); err != nil {
		return err
	} else if !isExists {
		errMsg := "EnvStore not found in path:" + currentEnvStoreFilePath
		return errors.New(errMsg)
	}

	if err := writeEnvMapToFile(currentEnvStoreFilePath, []EnvModel{}); err != nil {
		return err
	}

	return nil
}

func clearCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	if err := clearEnvs(); err != nil {
		log.Fatal("Failed to clear EnvStore:", err)
	}

	log.Info("EnvStore cleared")
}
