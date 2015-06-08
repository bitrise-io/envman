package main

import (
	//"fmt"
	"os"
	"os/exec"
)

type commandStruct struct {
	Command      string
	Argumentums  []string
	Environments envMap
}

func executeCmd(commandToRun commandStruct) error {
	cmd := exec.Command(commandToRun.Command, commandToRun.Argumentums...)
	//fmt.Println("cmd: ", cmd)

	cmdEnvs := []string{}
	for key, value := range commandToRun.Environments {
		cmdEnvs = append(cmdEnvs, key+"="+value)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), cmdEnvs...)

	return cmd.Run()
}
