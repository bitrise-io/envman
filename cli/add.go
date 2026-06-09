package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bitrise-io/envman/v2/envman"
	"github.com/bitrise-io/envman/v2/models"
	"github.com/bitrise-io/go-utils/pointers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const envVarLimitErrorKnowledgeBaseURL = "https://support.bitrise.io/en/articles/9676692-env-var-value-too-large-env-var-list-too-large"

func add(c *cli.Context) error {
	log.Debugln("[ENVMAN] Work path:", CurrentEnvStoreFilePath)

	key := c.String(KeyKey)
	expand := !c.Bool(NoExpandKey)
	replace := !c.Bool(AppendKey)
	skipIfEmpty := c.Bool(SkipIfEmptyKey)
	sensitive := c.Bool(SensitiveKey)

	var value string

	// read flag value
	if c.IsSet(ValueKey) {
		value = c.String(ValueKey)
		log.Debugf("adding flag value: (%s)", value)
	}

	// read flag file
	if value == "" && c.String(ValueFileKey) != "" {
		var err error
		if value, err = loadValueFromFile(c.String(ValueFileKey)); err != nil {
			log.Fatalf("[ENVMAN] Failed to read env var value from file: %s", err)
		}
		log.Debugf("adding file flag value: (%s)", value)
	}

	// read piped stdin value
	if value == "" {
		info, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("[ENVMAN] Failed to get file info for standard input: %s", err)
		}
		if info.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
			log.Debugf("adding from piped stdin")
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("[ENVMAN] Failed to read env var value from standard input: %s", err)
			}

			value = string(data)
			log.Debugf("stdin value: (%s)", value)
		}
	}

	if err := AddEnv(CurrentEnvStoreFilePath, key, value, expand, replace, skipIfEmpty, sensitive); err != nil {
		var envVarValueTooLargeErr EnvVarValueTooLargeError
		var envVarListTooLargeErr EnvVarListTooLargeError
		if errors.As(err, &envVarValueTooLargeErr) || errors.As(err, &envVarListTooLargeErr) {
			err = fmt.Errorf("%w.\nTo increase env var limits please visit: %s", err, envVarLimitErrorKnowledgeBaseURL)
		}
		log.Fatalf("[ENVMAN] Failed to add env var: %s", err)
	}

	log.Debugln("[ENVMAN] Env added")

	if err := logEnvs(CurrentEnvStoreFilePath); err != nil {
		log.Fatalf("[ENVMAN] Failed to print env var list: %s", err)
	}

	return nil
}

// AddEnv ...
func AddEnv(envStorePth string, key string, value string, expand, replace, skipIfEmpty, sensitive bool) error {
	if key == "" {
		return errors.New("key is not specified, required")
	}

	// Load envs, or create if not exist
	environments, err := ReadEnvsOrCreateEmptyList(envStorePth)
	if err != nil {
		return err
	}

	// Add or update envlist
	newEnv := models.EnvironmentItemModel{
		key: value,
		models.OptionsKey: models.EnvironmentItemOptionsModel{
			IsExpand:    pointers.NewBoolPtr(expand),
			SkipIfEmpty: pointers.NewBoolPtr(skipIfEmpty),
			IsSensitive: pointers.NewBoolPtr(sensitive),
		},
	}
	if err := newEnv.NormalizeValidateFillDefaults(); err != nil {
		return err
	}

	newEnvSlice, err := UpdateOrAddToEnvlist(environments, newEnv, replace)
	if err != nil {
		return err
	}

	// Validate the resulting env list as a whole, so the byte-limit overrides apply regardless
	// of where the override keys sit in the list.
	if err := validateEnvList(newEnvSlice); err != nil {
		return err
	}

	return WriteEnvMapToFile(envStorePth, newEnvSlice)
}

// validateEnvList enforces the byte limits over the whole env list. The two limit-override keys
// (EnvBytesLimitInKBEnvKey / EnvListBytesLimitInKBEnvKey) are pulled out of the list first and
// used in place of the configured limits, so their position in the list does not matter. The
// remaining env vars are then validated one by one against those limits.
func validateEnvList(envList []models.EnvironmentItemModel) error {
	configs, err := envman.GetConfigs()
	if err != nil {
		return err
	}

	envBytesLimitInKB := configs.EnvBytesLimitInKB
	envListBytesLimitInKB := configs.EnvListBytesLimitInKB

	// First pass: pull the limit overrides out of the list, wherever they are.
	for _, env := range envList {
		key, value, err := env.GetKeyValuePair()
		if err != nil {
			return err
		}
		switch key {
		case envman.EnvBytesLimitInKBEnvKey:
			if envBytesLimitInKB, err = envman.LimitFromEnvValue(value, key); err != nil {
				return err
			}
		case envman.EnvListBytesLimitInKBEnvKey:
			if envListBytesLimitInKB, err = envman.LimitFromEnvValue(value, key); err != nil {
				return err
			}
		}
	}

	// Second pass: validate every env var against the (possibly overridden) limits. The override
	// keys themselves are config, not payload, so they are excluded from the size checks.
	envListSizeInBytes := 0
	for _, env := range envList {
		key, value, err := env.GetKeyValuePair()
		if err != nil {
			return err
		}
		if key == envman.EnvBytesLimitInKBEnvKey || key == envman.EnvListBytesLimitInKBEnvKey {
			continue
		}

		valueSizeInBytes := len([]byte(value))
		if envBytesLimitInKB > 0 && valueSizeInBytes > envBytesLimitInKB*1024 {
			valueSizeInKB := (float64)(valueSizeInBytes) / 1024.0
			return NewEnvVarValueTooLargeError(key, valueSizeInKB, (float64)(envBytesLimitInKB))
		}

		envListSizeInBytes += valueSizeInBytes
	}

	if envListBytesLimitInKB > 0 && envListSizeInBytes > envListBytesLimitInKB*1024 {
		listSizeInKB := (float64)(envListSizeInBytes) / 1024.0
		return NewEnvVarListTooLargeError(listSizeInKB, (float64)(envListBytesLimitInKB))
	}

	return nil
}

func loadValueFromFile(pth string) (string, error) {
	buf, err := os.ReadFile(pth)
	if err != nil {
		return "", err
	}

	str := string(buf)
	return str, nil
}

func logEnvs(envStorePth string) error {
	environments, err := ReadEnvs(envStorePth)
	if err != nil {
		return err
	}

	if len(environments) == 0 {
		log.Info("[ENVMAN] Empty envstore")
	} else {
		for _, env := range environments {
			key, value, err := env.GetKeyValuePair()
			if err != nil {
				return err
			}

			opts, err := env.GetOptions()
			if err != nil {
				return err
			}

			envString := "- " + key + ": " + value
			log.Debugln(envString)
			if !*opts.IsExpand {
				expandString := "  " + "isExpand" + ": " + "false"
				log.Debugln(expandString)
			}
		}
	}

	return nil
}
