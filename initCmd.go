package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func initCmd(c *cli.Context) {
	fmt.Println("init")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	/*
		err := writeEnvMapToFile(envMapPath, envMap{})
		if err != nil {
			fmt.Println("Failed to clear envlist, err:%s", err)
		}
	*/
}
