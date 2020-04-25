package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/envman/envman"
	"github.com/urfave/cli"
)

func export(c *cli.Context) error {
	environments, err := envman.ReadEnvs(envman.CurrentEnvStoreFilePath)
	if err != nil {
		log.Fatalf("Failed to read envs, error: %s", err)
	}

	for _, item := range environments {
		key, value, err := item.GetKeyValuePair()
		if err != nil {
			continue
		}
		fmt.Printf("export %s=\"%s\"\n", key, value)
	}
	return nil
}
