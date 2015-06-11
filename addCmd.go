package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var addCmdLog *log.Entry = log.WithFields(log.Fields{"f": "addCmd.go"})

func addCmd(c *cli.Context) {
	key := c.String("key")
	value := c.String("value")
	if stdinValue != "" {
		value = stdinValue
	}
	value = strings.Replace(value, "\n", "", -1)

	// Validate input
	if key == "" {
		addCmdLog.Fatal("Invalid environment variable key")
	}
	if value == "" {
		addCmdLog.Fatal("Invalid environment variable value")
	}

	// Load envs, or create if not exist
	environments, err := loadEnvMapOrCreate()
	if err != nil {
		addCmdLog.Fatal("Failed to load envlist, err:", err)
	}

	// Add or update envlist
	newEnv := envMap{key: value}
	environments, err = updateOrAddToEnvlist(environments, newEnv)
	if err != nil {
		addCmdLog.Fatal("Failed to create store envlist, err:", err)
	}

	// Print new environment list
	addCmdLog.Info("Environment added, path:", currentEnvStorePath)
	printCmd(c)
}
