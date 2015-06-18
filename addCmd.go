package main

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func addCmd(c *cli.Context) {
	key := c.String(KEY_KEY)
	var value string

	if stdinValue != "" {
		value = stdinValue
	} else if c.String(VALUE_KEY) != "" {
		value = c.String(VALUE_KEY)
	} else if c.String(VALUE_FILE_KEY) != "" {
		v, err := loadValueFromFile(c.String(VALUE_FILE_KEY))
		if err != nil {
			log.Fatal("Failed to read file value: ", err)
		}
		value = v
	}

	// Validate input
	if key == "" {
		log.Fatal("Empty key")
	}
	if value == "" {
		log.Fatal("Empty value")
	}
	value = strings.Replace(value, "\n", "", -1)

	// Load envs, or create if not exist
	environments, err := loadEnvMapOrCreate()
	if err != nil {
		log.Fatal("Failed to load EnvStore:", err)
	}

	// Add or update envlist
	newEnv := envModel{key, value}
	environments, err = updateOrAddToEnvlist(environments, newEnv)
	if err != nil {
		log.Fatal("Failed to create EnvStore:", err)
	}

	// Print new environment list
	log.Info("Env added")
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
