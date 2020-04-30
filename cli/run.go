package cli

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/envman/envman"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	"github.com/urfave/cli"
)

type envController interface {
	Set(string, string) error
	Unset(string) error
	List() []string
	Expand(string) string
}

// CommandModel ...
type CommandModel struct {
	Command      string
	Argumentums  []string
	Environments []models.EnvironmentItemModel
}

func expandEnvsInString(inp string, ec envController) string {
	return ec.Expand(inp)
}

func commandEnvs(envs []models.EnvironmentItemModel, ec envController) ([]string, error) {
	for _, env := range envs {
		key, value, err := env.GetKeyValuePair()
		if err != nil {
			return []string{}, err
		}

		opts, err := env.GetOptions()
		if err != nil {
			return []string{}, err
		}

		if opts.Unset != nil && *opts.Unset {
			if err := ec.Unset(key); err != nil {
				return []string{}, fmt.Errorf("unset env (%s): %s", key, err)
			}
			continue
		}

		if *opts.SkipIfEmpty && value == "" {
			continue
		}

		var valueStr string
		if *opts.IsExpand {
			valueStr = expandEnvsInString(value, ec)
		} else {
			valueStr = value
		}

		if err := ec.Set(key, valueStr); err != nil {
			return []string{}, err
		}
	}
	return ec.List(), nil
}

func runCommandModel(cmdModel CommandModel) (int, error) {
	cmdEnvs, err := commandEnvs(cmdModel.Environments, EnvController{})
	if err != nil {
		return 1, err
	}

	return command.RunCommandWithEnvsAndReturnExitCode(cmdEnvs, cmdModel.Command, cmdModel.Argumentums...)
}

func run(c *cli.Context) error {
	log.Debug("[ENVMAN] - Work path:", envman.CurrentEnvStoreFilePath)

	if len(c.Args()) > 0 {
		doCmdEnvs, err := envman.ReadEnvs(envman.CurrentEnvStoreFilePath)
		if err != nil {
			log.Fatal("[ENVMAN] - Failed to load EnvStore:", err)
		}

		doCommand := c.Args()[0]

		doArgs := []string{}
		if len(c.Args()) > 1 {
			doArgs = c.Args()[1:]
		}

		cmdToExecute := CommandModel{
			Command:      doCommand,
			Environments: doCmdEnvs,
			Argumentums:  doArgs,
		}

		log.Debug("[ENVMAN] - Executing command:", cmdToExecute)

		if exit, err := runCommandModel(cmdToExecute); err != nil {
			log.Debug("[ENVMAN] - Failed to execute command:", err)
			if exit == 0 {
				log.Error("[ENVMAN] - Failed to execute command:", err)
				exit = 1
			}
			os.Exit(exit)
		}

		log.Debug("[ENVMAN] - Command executed")
	} else {
		log.Fatal("[ENVMAN] - No command specified")
	}

	return nil
}
