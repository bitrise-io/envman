package main

import (
	"fmt"
	"os"
	"os/exec"
)

type EnvironmentKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CommandModel struct {
	Command      string                `json:"command"`
	Argumentums  []string              `json:"argumentums"`
	Environments []EnvironmentKeyValue `json:"environments"`
}

func executeCmd(commandToRun CommandModel) error {

	cmd := exec.Command(commandToRun.Command, commandToRun.Argumentums...)
	fmt.Println("cmd: ", cmd)

	cmdEnvs := []string{}
	envLength := len(commandToRun.Environments)
	if envLength > 0 {
		cmdEnvs = make([]string, envLength, envLength)
		for idx, aEnvPair := range commandToRun.Environments {
			cmdEnvs[idx] = aEnvPair.Key + "=" + aEnvPair.Value
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), cmdEnvs...)

	return cmd.Run()
}
