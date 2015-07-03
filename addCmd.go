package main

import (
	"errors"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func addEnv(key string, value string, expand bool) error {
	// Validate input
	if key == "" {
		return errors.New("Key is not specified, required.")
	}

	// Load envs, or create if not exist
	environments, err := loadEnvMapOrCreate()
	if err != nil {
		return err
	}

	// Add or update envlist
	newEnv := envModel{key, value, expand}
	environments, err = updateOrAddToEnvlist(environments, newEnv)
	if err != nil {
		return err
	}

	return nil
}

func loadValueFromFile(pth string) (string, error) {
	buf, err := ioutil.ReadFile(pth)
	if err != nil {
		return "", err
	}

	str := string(buf)
	return str, nil
}

func addCmd(c *cli.Context) {
	log.Info("Work path:", currentEnvStoreFilePath)

	key := c.String(KEY_KEY)
	expand := isExpand(c.String(EXPAND_KEY))
	var value string

	if stdinValue != "" {
		value = stdinValue
	} else if c.IsSet(VALUE_KEY) {
		value = c.String(VALUE_KEY)
	} else if c.String(VALUE_FILE_KEY) != "" {
		v, err := loadValueFromFile(c.String(VALUE_FILE_KEY))
		if err != nil {
			log.Fatal("Failed to read file value: ", err)
		}
		value = v
	}

	err := addEnv(key, value, expand)
	if err != nil {
		log.Fatal("Failed to add env:", err)
	}

	log.Info("Env added")

	err = printEnvs()
	if err != nil {
		log.Fatal("Failed to print:", err)
	}
}
