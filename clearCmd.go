package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func clearCmd(c *cli.Context) {
	err := writeEnvMapToFile(envMapPath, envMap{})
	if err != nil {
		fmt.Println("Failed to clear envlist, err:%s", err)
	}
}
