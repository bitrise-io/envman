package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-io/envman/env"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CommandModel ...
type CommandModel struct {
	Command      string
	Argumentums  []string
	Environments []models.EnvironmentItemModel
}

func expandEnvsInString(inp string) string {
	return os.ExpandEnv(inp)
}

func commandEnvs(newEnvs []models.EnvironmentItemModel) ([]string, error) {
	result, err := env.GetDeclarationsSideEffects(newEnvs, &env.DefaultEnvironmentSource{})
	if err != nil {
		return nil, err
	}

	for _, command := range result.CommandHistory {
		if err := env.ExecuteCommand(command); err != nil {
			return nil, err
		}
	}

	return os.Environ(), nil
}

func runCommandModel(cmdModel CommandModel) (int, error) {
	cmdEnvs, err := commandEnvs(cmdModel.Environments)
	if err != nil {
		return 1, err
	}

	return command.RunCommandWithEnvsAndReturnExitCode(cmdEnvs, cmdModel.Command, cmdModel.Argumentums...)
}

func run(c *cli.Context) error {
	if len(c.Args()) == 0 {
		log.Fatal("[ENVMAN] - No command specified")
	}

	exitCode, err := RunCommand(CurrentEnvStoreFilePath, c.Args())
	if err != nil {
		log.Errorf("command failed: %s", err)
	}

	os.Exit(exitCode)

	return nil
}

func RunCommand(envStorePth string, args []string) (int, error) {
	if len(args) == 0 {
		return 1, fmt.Errorf("no command specified")
	}

	doCmdEnvs, err := ReadEnvs(envStorePth)
	if err != nil {
		return 1, fmt.Errorf("failed to load EnvStore: %s", err)
	}

	doCommand := args[0]

	doArgs := []string{}
	if len(args) > 1 {
		doArgs = args[1:]
	}

	cmdToExecute := CommandModel{
		Command:      doCommand,
		Environments: doCmdEnvs,
		Argumentums:  doArgs,
	}

	exit, err := runCommandModel(cmdToExecute)
	if err != nil && exit == 0 {
		exit = 1
	}
	return exit, err
}
