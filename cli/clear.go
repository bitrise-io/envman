package cli

import (
	"errors"

	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func clear(c *cli.Context) error {
	log.Debugln("[ENVMAN] - Work path:", CurrentEnvStoreFilePath)

	if err := ClearEnvs(CurrentEnvStoreFilePath); err != nil {
		log.Fatal("[ENVMAN] - Failed to clear EnvStore:", err)
	}

	log.Info("[ENVMAN] - EnvStore cleared")

	return nil
}

func ClearEnvs(envStorePth string) error {
	if isExists, err := pathutil.IsPathExists(envStorePth); err != nil {
		return err
	} else if !isExists {
		errMsg := "EnvStore not found in path:" + envStorePth
		return errors.New(errMsg)
	}

	return WriteEnvMapToFile(envStorePth, []models.EnvironmentItemModel{})
}
