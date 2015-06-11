package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

var initCmdLog *log.Entry = log.WithFields(log.Fields{"f": "initCmd.go"})

func initCmd(c *cli.Context) {
	var workDir string

	if currentEnvStorePath == "" {
		pth, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			initCmdLog.Fatal(err)
		}
		workDir = path.Join(pth, envStoreName)
	} else {
		workDir = currentEnvStorePath
	}

	workPath, err := initAtPath(workDir)
	if err != nil {
		initCmdLog.Fatal(err)
	}

	initCmdLog.Info("Envman initialized at:", workPath)
}

func initAtPath(pth string) (string, error) {
	ensuredPath, err := ensureEnvStorePath(pth)
	if err != nil {
		return "", err
	}

	exist, err := pathutil.IsPathExists(ensuredPath)
	if err != nil {
		return "", err
	}

	if exist == false {
		err = writeEnvMapToFile(ensuredPath, envMap{})
		if err != nil {
			return "", err
		}
	}

	return ensuredPath, nil
}

func ensureEnvStorePath(pth string) (string, error) {
	if path.Base(pth) == "." {
		return "", errors.New("No path sepcified")
	}

	_, file := path.Split(pth)
	if file == "" {
		return "", errors.New("Provided path is a directory not a file")
	}

	return pth, nil
}
