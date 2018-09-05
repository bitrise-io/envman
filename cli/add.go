package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/envman/envman"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/urfave/cli"
)

func envListSizeInBytes(envs []models.EnvironmentItemModel) (int, error) {
	valueSizeInBytes := 0
	for _, env := range envs {
		_, value, err := env.GetKeyValuePair()
		if err != nil {
			return 0, err
		}
		valueSizeInBytes += len([]byte(value))
	}
	return valueSizeInBytes, nil
}

func validateEnv(key, value string, envList []models.EnvironmentItemModel) (string, error) {
	if key == "" {
		return "", errors.New("Key is not specified, required")
	}

	configs, err := envman.GetConfigs()
	if err != nil {
		return "", err
	}

	valueSizeInBytes := len([]byte(value))
	if configs.EnvBytesLimitInKB > 0 {
		if valueSizeInBytes > configs.EnvBytesLimitInKB*1024 {
			valueSizeInKB := ((float64)(valueSizeInBytes)) / 1024.0
			log.Warnf("environment value (%s...) too large", value[0:100])
			log.Warnf("environment value size (%#v KB) - max allowed size: %#v KB", valueSizeInKB, (float64)(configs.EnvBytesLimitInKB))
			return "environment value too large - rejected", nil
		}
	}

	if configs.EnvListBytesLimitInKB > 0 {
		envListSizeInBytes, err := envListSizeInBytes(envList)
		if err != nil {
			return "", err
		}
		if envListSizeInBytes+valueSizeInBytes > configs.EnvListBytesLimitInKB*1024 {
			listSizeInKB := (float64)(envListSizeInBytes)/1024 + (float64)(valueSizeInBytes)/1024
			log.Warn("environment list too large")
			log.Warnf("environment list size (%#v KB) - max allowed size: %#v KB", listSizeInKB, (float64)(configs.EnvListBytesLimitInKB))
			return "", errors.New("environment list too large")
		}
	}
	return value, nil
}

func addEnv(key string, value string, expand, replace, skipIfEmpty bool) error {
	// Load envs, or create if not exist
	environments, err := envman.ReadEnvsOrCreateEmptyList()
	if err != nil {
		return err
	}

	// Validate input
	validatedValue, err := validateEnv(key, value, environments)
	if err != nil {
		return err
	}
	value = validatedValue

	// Add or update envlist
	newEnv := models.EnvironmentItemModel{
		key: value,
		models.OptionsKey: models.EnvironmentItemOptionsModel{
			IsExpand:    pointers.NewBoolPtr(expand),
			SkipIfEmpty: pointers.NewBoolPtr(skipIfEmpty),
		},
	}
	if err := newEnv.NormalizeValidateFillDefaults(); err != nil {
		return err
	}

	newEnvSlice, err := envman.UpdateOrAddToEnvlist(environments, newEnv, replace)
	if err != nil {
		return err
	}

	return envman.WriteEnvMapToFile(envman.CurrentEnvStoreFilePath, newEnvSlice)
}

func loadValueFromFile(pth string) (string, error) {
	buf, err := ioutil.ReadFile(pth)
	if err != nil {
		return "", err
	}

	str := string(buf)
	return str, nil
}

func logEnvs() error {
	environments, err := envman.ReadEnvs(envman.CurrentEnvStoreFilePath)
	if err != nil {
		return err
	}

	if len(environments) == 0 {
		log.Info("[ENVMAN] - Empty envstore")
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

func readIfNamedPipe(f *os.File) (string, bool, error) {
	info, err := f.Stat()
	if err != nil {
		return "", false, fmt.Errorf("failed to get stdin FileInfo: %s", err)
	}

	var stdinValue string
	var isNamedPipe bool
	switch {
	case info.Mode()&os.ModeNamedPipe == os.ModeNamedPipe:
		log.Debugln("reading a named pipe")
		if info.Size() > 0 {
			var err error
			stdinValue, err = readAllByRunes(f)
			if err != nil {
				return "", false, fmt.Errorf("failed to read: %s", err)
			}
			isNamedPipe = true
		}
	case info.Mode()&os.ModeCharDevice == os.ModeCharDevice:
		log.Debugln("reading a char device")
	default:
		log.Debugf("reading: %s", info.Mode())
	}
	log.Debugf("stdin has %d bytes", info.Size())
	log.Debugf("stdinValue: <%s>", stdinValue)
	return stdinValue, isNamedPipe, nil
}

func readAllByRunes(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	var output []rune
	for {
		c, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to read rune: %s", err)
		}
		output = append(output, c)
	}
	return string(output), nil
}

func add(c *cli.Context) error {
	log.Debugln("[ENVMAN] - Work path:", envman.CurrentEnvStoreFilePath)

	key := c.String(KeyKey)
	expand := !c.Bool(NoExpandKey)
	replace := !c.Bool(AppendKey)
	skipIfEmpty := c.Bool(SkipIfEmptyKey)

	var value string
	if stdinValue, _, err := readIfNamedPipe(os.Stdin); err != nil {
		log.Fatal("Failed to read stdin:", err)
	} else if stdinValue != "" {
		value = stdinValue
	} else if c.IsSet(ValueKey) {
		value = c.String(ValueKey)
	} else if c.String(ValueFileKey) != "" {
		if v, err := loadValueFromFile(c.String(ValueFileKey)); err != nil {
			log.Fatal("[ENVMAN] - Failed to read file value: ", err)
		} else {
			value = v
		}
	}

	if err := addEnv(key, value, expand, replace, skipIfEmpty); err != nil {
		log.Fatal("[ENVMAN] - Failed to add env:", err)
	}

	log.Debugln("[ENVMAN] - Env added")

	if err := logEnvs(); err != nil {
		log.Fatal("[ENVMAN] - Failed to print:", err)
	}

	return nil
}
