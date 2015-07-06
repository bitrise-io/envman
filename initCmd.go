package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

func initAtPath(pth string) error {
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return err
	} else if exist == false {
		if err := writeEnvMapToFile(pth, []EnvModel{}); err != nil {
			return err
		}
	}
	return nil
}

func initCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	if err := initAtPath(currentEnvStoreFilePath); err != nil {
		log.Fatal(err)
	}

	log.Info("Envman initialized")
}
