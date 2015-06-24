package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	ENVMAN_ENVSTORE_PATH_KEY string = "ENVMAN_ENVSTORE_PATH"
	envStoreName             string = ".envstore.yml"
)

var (
	stdinValue              string
	currentEnvStoreFilePath string
)

func isPipedData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}

// Run the Envman CLI.
func run() {
	log.SetLevel(log.DebugLevel)

	// Read piped data
	if isPipedData() {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err == nil {
			if len(bytes) > 0 {
				stdinValue = string(bytes)
			}
		} else {
			log.Error("Failed to read stdin:", err)
		}
	}

	// Parse cl
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Environment variable manager"
	app.Version = VERSION

	app.Author = ""
	app.Email = ""

	app.Before = func(c *cli.Context) error {
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatal(err.Error())
		}
		log.SetLevel(level)

		// Befor parsing cli, and running command
		// we need to decide wich path will be used by envman
		flagPath := c.String(PATH_KEY)
		if flagPath == "" {
			if os.Getenv(ENVMAN_ENVSTORE_PATH_KEY) != "" {
				currentEnvStoreFilePath = os.Getenv(ENVMAN_ENVSTORE_PATH_KEY)
			} else {
				currentPath, err := envStoreInCurrentPath()
				if err != nil {
					log.Fatal("Failed to set envman work path in current dir:", err)
				}
				currentEnvStoreFilePath = currentPath
			}

			return nil
		}

		if err := validatePath(flagPath); err != nil {
			log.Fatal("Failed to set envman work path:", err)
		}

		currentEnvStoreFilePath = flagPath

		return nil
	}

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Envman finished:", err)
	}
}

/*
Check if current path contains .envstore.yml
Output :
	@string: - current envman work path
	@error: - error
*/
func envStoreInCurrentPath() (string, error) {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return currentDir, err
	}

	currentPath := path.Join(currentDir, envStoreName)
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
		return errors.New("No path sepcified")
	}
	_, file := path.Split(pth)
	if file == "" {
		return errors.New("EnvStore not found in path:" + pth)
	}
	return nil
}
