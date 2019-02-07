package cli

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/envman/envman"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	"github.com/urfave/cli"
)

// CommandModel ...
type CommandModel struct {
	Command     string
	Argumentums []string
	EnvStore    models.EnvsSerializeModel
}

func expandEnvsInString(inp string) string {
	return os.ExpandEnv(inp)
}

func commandEnvs(envs models.EnvsSerializeModel) ([]string, error) {
	for _, env := range envs.Envs {
		key, value, err := env.GetKeyValuePair()
		if err != nil {
			return []string{}, err
		}

		opts, err := env.GetOptions()
		if err != nil {
			return []string{}, err
		}

		if *opts.SkipIfEmpty && value == "" {
			continue
		}

		var valueStr string
		if *opts.IsExpand {
			valueStr = expandEnvsInString(value)
		} else {
			valueStr = value
		}

		if err := os.Setenv(key, valueStr); err != nil {
			return []string{}, err
		}
	}

	for _, key := range envs.Unsets {
		os.Unsetenv(key)
	}

	return os.Environ(), nil
}

func runCommandModel(cmdModel CommandModel) (int, error) {
	cmdEnvs, err := commandEnvs(cmdModel.EnvStore)
	if err != nil {
		return 1, err
	}

	return command.RunCommandWithEnvsAndReturnExitCode(cmdEnvs, cmdModel.Command, cmdModel.Argumentums...)
}

func run(c *cli.Context) error {
	log.Debug("[ENVMAN] - Work path:", envman.CurrentEnvStoreFilePath)

	if len(c.Args()) > 0 {
		doEnvstore, err := envman.ReadEnvs(envman.CurrentEnvStoreFilePath)
		if err != nil {
			log.Fatal("[ENVMAN] - Failed to load EnvStore:", err)
		}

		doCommand := c.Args()[0]

		doArgs := []string{}
		if len(c.Args()) > 1 {
			doArgs = c.Args()[1:]
		}

		cmdToExecute := CommandModel{
			Command:     doCommand,
			EnvStore:    doEnvstore,
			Argumentums: doArgs,
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
