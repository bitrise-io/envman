package main

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var addCmdLog *log.Entry = log.WithFields(log.Fields{"f": "addCmd.go"})

func addCmd(c *cli.Context) {
	key := c.String(KEY_KEY)
	expand := isExpand(c.String(EXPAND_KEY))
	var value string

	if stdinValue != "" {
		value = stdinValue
	} else if c.String(VALUE_KEY) != "" {
		value = c.String(VALUE_KEY)
	} else if c.String(VALUE_FILE_KEY) != "" {
		v, err := loadValueFromFile(c.String(VALUE_FILE_KEY))
		if err != nil {
			addCmdLog.Fatal("Failed to read file value, err: ", err)
		}
		value = v
	}

	// Validate input
	if key == "" {
		addCmdLog.Fatal("Invalid environment variable key")
	}
	if value == "" {
		addCmdLog.Fatal("Invalid environment variable value")
	}
	value = strings.Replace(value, "\n", "", -1)

	// Load envs, or create if not exist
	environments, err := loadEnvMapOrCreate()
	if err != nil {
		addCmdLog.Fatal("Failed to load envlist, err:", err)
	}

	// Add or update envlist
	newEnv := envModel{key, value, expand}
	addCmdLog.Info("envs, newEnv: ", environments, newEnv)
	environments, err = updateOrAddToEnvlist(environments, newEnv)
	if err != nil {
		addCmdLog.Fatal("Failed to create store envlist, err:", err)
	}

	// Print new environment list
	addCmdLog.Info("Environment added, path:", currentEnvStoreFilePath)
	printCmd(c)
}

func loadValueFromFile(pth string) (string, error) {
	buf, err := ioutil.ReadFile(pth)
	if err != nil {
		return "", err
	}

	str := string(buf)
	return str, nil
}
