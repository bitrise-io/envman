package envman

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	envmanConfigFileName         = "configs.json"
	defaultEnvBytesLimitInKB     = 256
	defaultEnvListBytesLimitInKB = 256

	// EnvBytesLimitInKBEnvKey and EnvListBytesLimitInKBEnvKey are the env keys to override the
	// byte limits. Env-based overrides take precedence over the config file, so the limits stay
	// overrideable when the config file is not writable/reachable, e.g. while running a script step.
	// They are read both from the process environment and from envman's own env list (envstore).
	EnvBytesLimitInKBEnvKey     = "ENVMAN_ENV_BYTES_LIMIT_IN_KB"
	EnvListBytesLimitInKBEnvKey = "ENVMAN_ENV_LIST_BYTES_LIMIT_IN_KB"
)

// ConfigsModel ...
type ConfigsModel struct {
	EnvBytesLimitInKB     int `json:"env_bytes_limit_in_kb,omitempty"`
	EnvListBytesLimitInKB int `json:"env_list_bytes_limit_in_kb,omitempty"`
}

func getEnvmanConfigsDirPath() string {
	return path.Join(pathutil.UserHomeDir(), ".envman")
}

func getEnvmanConfigsFilePath() string {
	return path.Join(getEnvmanConfigsDirPath(), envmanConfigFileName)
}

func ensureEnvmanConfigDirExists() error {
	confDirPth := getEnvmanConfigsDirPath()
	isExists, err := pathutil.IsDirExists(confDirPth)
	if !isExists || err != nil {
		if err := os.MkdirAll(confDirPth, 0777); err != nil {
			return err
		}
	}
	return nil
}

func createDefaultConfigsModel() ConfigsModel {
	return ConfigsModel{
		EnvBytesLimitInKB:     defaultEnvBytesLimitInKB,
		EnvListBytesLimitInKB: defaultEnvListBytesLimitInKB,
	}
}

// GetConfigs ...
func GetConfigs() (ConfigsModel, error) {
	configPth := getEnvmanConfigsFilePath()
	configs := createDefaultConfigsModel()

	isExist, err := pathutil.IsPathExists(configPth)
	if err != nil {
		return ConfigsModel{}, err
	}

	if isExist {
		bytes, err := fileutil.ReadBytesFromFile(configPth)
		if err != nil {
			return ConfigsModel{}, err
		}

		type ConfigsFileMode struct {
			EnvBytesLimitInKB     *int `json:"env_bytes_limit_in_kb,omitempty"`
			EnvListBytesLimitInKB *int `json:"env_list_bytes_limit_in_kb,omitempty"`
		}

		var userConfigs ConfigsFileMode
		if err := json.Unmarshal(bytes, &userConfigs); err != nil {
			return ConfigsModel{}, err
		}

		if userConfigs.EnvBytesLimitInKB != nil {
			configs.EnvBytesLimitInKB = *userConfigs.EnvBytesLimitInKB
		}
		if userConfigs.EnvListBytesLimitInKB != nil {
			configs.EnvListBytesLimitInKB = *userConfigs.EnvListBytesLimitInKB
		}
	}

	// Process environment variables take precedence over the config file, so the limits stay
	// overrideable even when the config file is not reachable (e.g. inside a script step).
	// Overrides coming from envman's own env list (see OverrideConfigsWithEnvs) take precedence
	// over these, as that is where the override usually lives during a build.
	if err := OverrideConfigsWithEnvs(&configs, os.LookupEnv); err != nil {
		return ConfigsModel{}, err
	}

	return configs, nil
}

// OverrideConfigsWithEnvs overrides the byte limits from env values resolved by lookup.
// lookup mirrors os.LookupEnv: it returns the value and whether the key is set. This lets the
// limits be overridden both from the process environment and from envman's own env list (envstore),
// which is important when envman cannot reach its config file, e.g. while running a script step.
func OverrideConfigsWithEnvs(configs *ConfigsModel, lookup func(key string) (string, bool)) error {
	if val, ok := lookup(EnvBytesLimitInKBEnvKey); ok {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid value (%s) for %s: %s", val, EnvBytesLimitInKBEnvKey, err)
		}
		configs.EnvBytesLimitInKB = limit
	}
	if val, ok := lookup(EnvListBytesLimitInKBEnvKey); ok {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid value (%s) for %s: %s", val, EnvListBytesLimitInKBEnvKey, err)
		}
		configs.EnvListBytesLimitInKB = limit
	}
	return nil
}

// saveConfigs ...
//  only used for unit testing at the moment
func saveConfigs(configModel ConfigsModel) error {
	if err := ensureEnvmanConfigDirExists(); err != nil {
		return err
	}

	bytes, err := json.Marshal(configModel)
	if err != nil {
		return err
	}
	configsPth := getEnvmanConfigsFilePath()
	return fileutil.WriteBytesToFile(configsPth, bytes)
}
