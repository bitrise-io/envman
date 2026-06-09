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
	// byte limits. When present in envman's env list (envstore), their values override the limits
	// from the config file during validation, so the limits stay overrideable when the config file
	// is not writable/reachable, e.g. while running a script step. The position of these keys in
	// the env list does not matter; they are pulled out before the rest of the list is validated.
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

	return configs, nil
}

// LimitFromEnvValue parses a byte-limit override value (in KB). key is only used for error context.
func LimitFromEnvValue(value, key string) (int, error) {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid value (%s) for %s: %s", value, key, err)
	}
	return limit, nil
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
