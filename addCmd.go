package main

import (
	"log"
	"strings"

	"github.com/codegangsta/cli"
)

func addCmd(c *cli.Context) {
	key := c.String("key")
	value := c.String("value")
	if stdinValue != "" {
		value = stdinValue
	}
	value = strings.Replace(value, "\n", "", -1)

	// Validate input
	if key == "" {
		log.Fatalln("Invalid environment variable key")
	}
	if value == "" {
		log.Fatalln("Invalid environment variable value")
	}

	// Load envs, or create if not exist
	environments, err := loadEnvMapOrCreate()
	if err != nil {
		log.Fatalln("Failed to load envlist, err:", err)
	}

	// Add or update envlist
	newEnv := envMap{key: value}
	environments, err = updateOrAddToEnvlist(environments, newEnv)
	if err != nil {
		log.Fatalln("Failed to create store envlist, err:", err)
	}

	// Print new environment list
	printCmd(c)
}
