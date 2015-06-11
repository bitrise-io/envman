package main

import (
	"os"
	"os/exec"

	//log "github.com/Sirupsen/logrus"
)

type commandModel struct {
	Command      string
	Argumentums  []string
	Environments envMap
}

func executeCmd(commandToRun commandModel) error {
	cmdEnvs := []string{}
	for key, value := range commandToRun.Environments {
		cmdEnvs = append(cmdEnvs, key+"="+value)
	}

	cmd := exec.Command(commandToRun.Command, commandToRun.Argumentums...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), cmdEnvs...)

	return cmd.Run()
}
