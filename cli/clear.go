package cli

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/envman/envman"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
	"gopkg.in/urfave/cli.v2"
)

func clearEnvs() error {
	if isExists, err := pathutil.IsPathExists(envman.CurrentEnvStoreFilePath); err != nil {
		return err
	} else if !isExists {
		errMsg := "EnvStore not found in path:" + envman.CurrentEnvStoreFilePath
		return errors.New(errMsg)
	}

	if err := envman.WriteEnvMapToFile(envman.CurrentEnvStoreFilePath, []models.EnvironmentItemModel{}); err != nil {
		return err
	}

	return nil
}

func clear(c *cli.Context) {
	log.Debugln("[ENVMAN] - Work path:", envman.CurrentEnvStoreFilePath)

	if err := clearEnvs(); err != nil {
		log.Fatal("[ENVMAN] - Failed to clear EnvStore:", err)
	}

	log.Info("[ENVMAN] - EnvStore cleared")
}
