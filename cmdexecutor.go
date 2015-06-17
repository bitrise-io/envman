package main

import (
	"os"
	"os/exec"

	//log "github.com/Sirupsen/logrus"
)

//var cmdexecutorLog *log.Entry = log.WithFields(log.Fields{"f": "cmdexecutor.go"})

type commandModel struct {
	Command      string
	Argumentums  []string
	Environments []envModel
}

func ExpandEnvsInString(inp string) string {
	return os.ExpandEnv(inp)
}

func executeCmd(commandToRun commandModel) error {
	cmdEnvs := []string{}
	for _, eModel := range commandToRun.Environments {
		var value string
		key := eModel.Key
		if eModel.IsExpand {
			value = ExpandEnvsInString(eModel.Value)
		} else {
			value = eModel.Value
		}

		os.Setenv(key, value)
		cmdEnvs = append(cmdEnvs, key+"="+value)
	}

	cmd := exec.Command(commandToRun.Command, commandToRun.Argumentums...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), cmdEnvs...)

	return cmd.Run()
}
