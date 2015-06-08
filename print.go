package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

func print(c *cli.Context) {
	environments, err := loadEnvMap()
	if err != nil {
		log.Fatalln("Failed to export environment variable list, err:", err)
	}
	if len(environments) == 0 {
		fmt.Println("Empty environment variable list")
	} else {
		fmt.Println(environments)
	}
}
