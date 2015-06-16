package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh/terminal"
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

const envStoreName string = ".envstore.yml"

var (
	cliLog                  *log.Entry = log.WithFields(log.Fields{"f": "cli.go"})
	stdinValue              string
	currentEnvStoreFilePath string // !!! keep in mind this should be like {SOME_DIR/envstore.yml}
)

// Run the Envman CLI.
func run() {
	// Read piped data
	if !terminal.IsTerminal(0) {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			cliLog.Fatal("Failed to read stdin, err:", err)
		}
		stdinValue = string(bytes)
	}

	// Parse cl
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Environment variable manager"
	app.Version = VERSION

	app.Author = ""
	app.Email = ""

	app.Before = func(c *cli.Context) error {
		// Befor parsing cli, and running command
		// we need to decide wich path will be used by envman
		flagPath := c.String("path")
		if flagPath == "" {
			if os.Getenv("ENVMAN_ENVSTORE_PATH") != "" {
				currentEnvStoreFilePath = os.Getenv("ENVMAN_ENVSTORE_PATH")
			} else {
				currentPath, err := ensureEnvStoreInCurrentPath()
				if err != nil {
					cliLog.Error(err)
				}
				currentEnvStoreFilePath = currentPath
			}
			cliLog.Info("Envman work path : %v", currentEnvStoreFilePath)
			return nil
		}

		if err := validatePath(flagPath); err != nil {
			cliLog.Fatal("Failed to set envman work path to: %s, err: %s", flagPath, err)
			return nil
		}

		currentEnvStoreFilePath = flagPath
		cliLog.Info("Envman work path : %v", currentEnvStoreFilePath)
		return nil
	}

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		cliLog.Fatal(err)
	}
}

/*
Check if current path contains .envstore.yml
Output :
	@string: - current envman work path
	@error: - error
*/
func ensureEnvStoreInCurrentPath() (string, error) {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return currentDir, err
	}

	currentPath := path.Join(currentDir, envStoreName)
	exist, err := pathutil.IsPathExists(currentPath)
	if err != nil {
		return currentPath, err
	}
	if !exist {
		err = errors.New(".envstore.yml dos not exist in current path: " + currentPath)
		return currentPath, err
	}

	return currentPath, nil
}

/*
Check if path is valid (i.e is not empty, and not a directory)
Input:
	@pth string - the path to validate
Output:
	@error - path is empty or not valid envstore file path
*/
func validatePath(pth string) error {
	if pth == "" {
		return errors.New("No path sepcified, should be like {SOME_DIR/.envstore.yml}")
	}
	_, file := path.Split(pth)
	if file == "" {
		return errors.New("EnvStore not found")
	}
	return nil
}
