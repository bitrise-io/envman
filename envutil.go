package main

import (
	"errors"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil"
	"gopkg.in/yaml.v2"
)

/*
	File storage methods
*/
func loadEnvMap() ([]envModel, error) {
	envsYML, err := readEnvMapFromFile(currentEnvStoreFilePath)
	if err != nil {
		return []envModel{}, err
	}

	return convertToEnvModelArray(envsYML), nil
}

func loadEnvMapOrCreate() ([]envModel, error) {
	envModels, err := loadEnvMap()
	if err != nil {
		if err.Error() == "No environment variable list found" {
			_, err = initAtPath(currentEnvStoreFilePath)
			return []envModel{}, err
		}
		return []envModel{}, err
	}
	return envModels, nil
}

func updateOrAddToEnvlist(envs []envModel, env envModel) ([]envModel, error) {
	var newEnvs []envModel
	exist := false

	for _, eModel := range envs {
		if eModel.Key == env.Key {
			exist = true
			newEnvs = append(newEnvs, env)
		} else {
			newEnvs = append(newEnvs, eModel)
		}
	}

	if exist == false {
		newEnvs = append(newEnvs, env)
	}

	err := writeEnvMapToFile(currentEnvStoreFilePath, newEnvs)
	if err != nil {
		return []envModel{}, err
	}

	return newEnvs, nil
}

func readEnvMapFromFile(pth string) (envsYMLModel, error) {
	isExists, err := pathutil.IsPathExists(pth)
	if err != nil {
		return envsYMLModel{}, err
	}
	if isExists == false {
		return envsYMLModel{}, errors.New("No environment variable list found")
	}

	bytes, err := ioutil.ReadFile(pth)
	if err != nil {
		return envsYMLModel{}, err
	}

	var envsModel envsYMLModel
	err = yaml.Unmarshal(bytes, &envsModel)
	if err != nil {
		return envsYMLModel{}, err
	}

	return envsModel, nil
}

func generateFormattedYMLForEnvModels(envs []envModel) ([]byte, error) {
	envYML := convertToEnvsYMLModel(envs)

	bytes, err := yaml.Marshal(envYML)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func writeEnvMapToFile(pth string, envs []envModel) error {
	if pth == "" {
		return errors.New("No path provided")
	}

	file, err := os.Create(pth)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalln("Failed to close file:", err)
		}
	}()

	jsonContBytes, err := generateFormattedYMLForEnvModels(envs)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonContBytes)
	if err != nil {
		return err
	}

	return nil
}
