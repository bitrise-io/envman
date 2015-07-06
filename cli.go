package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	defaultEnvStoreName string = ".envstore.yml"
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

func envStorePathInCurrentDir() (string, error) {
	return filepath.Abs(path.Join("./", defaultEnvStoreName))
}

func before(c *cli.Context) error {
	level, err := log.ParseLevel(c.String("log-level"))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.SetLevel(level)

	// Befor parsing cli, and running command
	// we need to decide wich path will be used by envman
	currentEnvStoreFilePath = c.String(PATH_KEY)
	if currentEnvStoreFilePath == "" {
		currentEnvStoreFilePath, err = envStorePathInCurrentDir()
		if err != nil {
			log.Fatal("Failed to set envman work path in current dir:", err)
		}
		return nil
	}

	return nil
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

	app.Before = before

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Envman finished:", err)
	}
}
