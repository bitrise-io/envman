package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

func initAtPath(pth string) error {
	exist, err := pathutil.IsPathExists(pth)
	if err != nil {
		return err
	}

	if exist == false {
		err = writeEnvMapToFile(pth, []EnvModel{})
		if err != nil {
			return err
		}
	}

	return nil
}

func initCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	err := initAtPath(currentEnvStoreFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Envman initialized")
}
