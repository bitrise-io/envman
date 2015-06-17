package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

func initCmd(c *cli.Context) {
	_, err := initAtPath(currentEnvStoreFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Envman initialized")
}

/*
	Initialize envman in specified path. Creates empty envstore if does not exist.
*/
func initAtPath(pth string) (string, error) {
	if err := validatePath(pth); err != nil {
		return "", err
	}

	exist, err := pathutil.IsPathExists(pth)
	if err != nil {
		return "", err
	}

	if exist == false {
		err = writeEnvMapToFile(pth, []envModel{})
		if err != nil {
			return "", err
		}
	}

	return pth, nil
}
