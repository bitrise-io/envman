package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh/terminal"
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"github.com/codegangsta/cli"
)

var cliLog *log.Entry = log.WithFields(log.Fields{"f": "cli.go"})

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
	app.Usage = "Environment varaibale manager"
	app.Version = VERSION

	app.Author = ""
	app.Email = ""

	app.Before = func(c *cli.Context) error {
		var err error
		if c.String("path") == "" {
			currentEnvStorePath, err = envStorePath()
			return err
		}

		ensuredPath, err := ensureEnvStorePath(c.String("path"))
		if err != nil {
			cliLog.Error(err)
			currentEnvStorePath, err = envStorePath()
			return err
		}

		currentEnvStorePath = ensuredPath
		return nil
	}
	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		cliLog.Fatal(err)
	}
}

func envStorePath() (string, error) {
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		cliLog.Fatal(err)
	}
	workPath := path.Join(workDir, envStoreName)

	paths := []string{workPath, defaultEnvStorePath}

	for _, path := range paths {
		exist, err := pathutil.IsPathExists(path)
		if err != nil || !exist {
			continue
		}
		return path, nil
	}

	err = createDeafultEnvmanDir()
	if err != nil {
		cliLog.Fatal(err)
	}

	return defaultEnvStorePath, nil
}
