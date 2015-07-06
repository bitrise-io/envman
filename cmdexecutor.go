package main

import (
	"os"
	"os/exec"
)

type commandModel struct {
	Command      string
	Argumentums  []string
	Environments []EnvModel
}

func expandEnvsInString(inp string) string {
	return os.ExpandEnv(inp)
}

func commandEnvs(envs []EnvModel) ([]string, error) {
	cmdEnvs := []string{}

	for _, eModel := range envs {
		var value string

		if eModel.IsExpand {
			value = expandEnvsInString(eModel.Value)
		} else {
			value = eModel.Value
		}

		if err := os.Setenv(eModel.Key, eModel.Value); err != nil {
			return []string{}, err
		}
		cmdEnvs = append(cmdEnvs, eModel.Key+"="+value)
	}

	return os.Environ(), nil
}

func executeCmd(commandToRun commandModel) error {
	cmdEnvs, err := commandEnvs(commandToRun.Environments)
	if err != nil {
		return err
	}

	cmd := exec.Command(commandToRun.Command, commandToRun.Argumentums...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = cmdEnvs

	return cmd.Run()
}
