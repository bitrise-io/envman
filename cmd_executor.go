package main

import (
	"os"
	"os/exec"
	"fmt"
	//"github.com/codeskyblue/go-sh"
)

type EnvironmentKeyValue struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type CommandModel struct {
	Command string `json:"command"`
	Argumentums []string `json:"argumentums"`
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

	//fmt.Println(cmd.Env)

	return cmd.Run()
	
	/*
	session := sh.NewSession()
	
	envLength := len(commandToRun.Environments)
	if envLength > 0 {
		for _, aEnvPair := range commandToRun.Environments {
			session.SetEnv(aEnvPair.Key, aEnvPair.Value)
		}
	}

	fmt.Println(session.Env)

	session.Command(commandToRun.Command, commandToRun.Argumentums).Run()
	*/

	return nil
}