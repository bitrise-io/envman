package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	//log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"gopkg.in/yaml.v2"
)

type envMap map[string]string

type environmentsModel struct {
	Environments envMap `yml:"environments"`
}

const (
	envStoreName string = ".envstore.yml"
)

var (
	defaultEnvmanDir    string = pathutil.UserHomeDir() + "/.envman/"
	defaultEnvStorePath string = defaultEnvmanDir + envStoreName
	currentEnvStorePath string
	stdinValue          string
)

//var envutilLog *log.Entry = log.WithFields(log.Fields{"f": "envutil.go"})

func createDeafultEnvmanDir() error {
	dir, _ := path.Split(defaultEnvStorePath)
	exist, err := pathutil.IsPathExists(dir)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func createEnvmanDir() error {
	dir, _ := path.Split(currentEnvStorePath)
	exist, err := pathutil.IsPathExists(dir)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func loadEnvMap() (envMap, error) {
	environments, err := readEnvMapFromFile(currentEnvStorePath)
	if err != nil {
		return envMap{}, err
	}

	return environments, nil
}

func loadEnvMapOrCreate() (envMap, error) {
	environments, err := loadEnvMap()
	if err != nil {
		if err.Error() == "No environment variable list found" {
			err = createEnvmanDir()
			return envMap{}, err
		}
		return envMap{}, err
	}
	return environments, nil
}

func updateOrAddToEnvlist(environments envMap, newEnv envMap) (envMap, error) {
	newEnvironments := make(envMap)
	for key, value := range environments {
		newEnvironments[key] = value
	}
	for key, value := range newEnv {
		newEnvironments[key] = value
	}

	err := writeEnvMapToFile(currentEnvStorePath, newEnvironments)
	if err != nil {
		return envMap{}, err
	}

	return newEnvironments, nil
}

func readEnvMapFromFile(path string) (envMap, error) {
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		return envMap{}, err
	}
	if isExists == false {
		return envMap{}, errors.New("No environment variable list found")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return envMap{}, err
	}

	var envs environmentsModel
	err = yaml.Unmarshal(bytes, &envs)
	if err != nil {
		return envMap{}, err
	}

	return envs.Environments, nil
}

func generateFormattedYMLForEnvMap(environments envMap) ([]byte, error) {
	var envs environmentsModel
	envs.Environments = environments

	bytes, err := yaml.Marshal(envs)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func writeEnvMapToFile(path string, environments envMap) error {
	if path == "" {
		return errors.New("No path provided")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonContBytes, err := generateFormattedYMLForEnvMap(environments)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonContBytes)
	if err != nil {
		return err
	}

	return nil
}
