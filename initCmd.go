package main

import (
	//"os"
	//"path"
	//"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

var initCmdLog *log.Entry = log.WithFields(log.Fields{"f": "initCmd.go"})

func initCmd(c *cli.Context) {
	workPath, err := initAtPath(currentEnvStoreFilePath)
	if err != nil {
		initCmdLog.Fatal(err)
	}

	initCmdLog.Info("Envman initialized at:", workPath)
}

func initAtPath(pth string) (string, error) {
	if err := validatePath(pth); err != nil {
		return "", err
	}

	exist, err := pathutil.IsPathExists(pth)
	if err != nil {
		return "", err
	}

	if exist == false {
		err = writeEnvMapToFile(pth, envMap{})
		if err != nil {
			return "", err
		}
	}

	return pth, nil
}
