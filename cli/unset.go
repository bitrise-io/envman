package cli

import (
	"github.com/bitrise-io/envman/envman"
	"github.com/urfave/cli"
)

func unset(c *cli.Context) error {
	key := c.String(KeyKey)

	envstore, err := envman.ReadEnvsOrCreateEmptyList()
	if err != nil {
		return err
	}

	envstore.Unsets = append(envstore.Unsets, key)

	return envman.WriteEnvMapToFile(envman.CurrentEnvStoreFilePath, envstore.Envs, envstore.Unsets)
}
