package main

import (
	"os"
	"os/exec"
)

type commandModel struct {
	Command      string
	Argumentums  []string
	Environments []envModel
}

func expandEnvsInString(inp string) string {
	return os.ExpandEnv(inp)
}

func commandEnvs(envs []envModel) ([]string, error) {
	cmdEnvs := []string{}

	// Exporting envs to bash_profile is required for expanding envs
	for _, eModel := range envs {
		err := os.Setenv(eModel.Key, eModel.Value)
		if err != nil {
			return cmdEnvs, err
		}
	}

	for _, eModel := range envs {
		var value string
		key := eModel.Key
		if eModel.IsExpand {
			value = expandEnvsInString(eModel.Value)
		} else {
			value = eModel.Value
		}

		cmdEnvs = append(cmdEnvs, key+"="+value)
	}

	return append(os.Environ(), cmdEnvs...), nil
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
