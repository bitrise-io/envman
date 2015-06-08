package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/bitrise-io/go-pathutil"
	"gopkg.in/yaml.v2"
)

type envMap map[string]string

type environmentsStruct struct {
	Environments envMap `yml:"environments"`
}

func readEnvMapFromFile(path string) (envMap, error) {
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		return envMap{}, err
	}
	if isExists == false {
		return envMap{}, errors.New("No environemt variable list found")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return envMap{}, err
	}

	var envs environmentsStruct
	err = yaml.Unmarshal(bytes, &envs)
	if err != nil {
		return envMap{}, err
	}

	return envs.Environments, nil
}

func generateFormattedYMLForEnvMap(environments envMap) ([]byte, error) {
	var envs environmentsStruct
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
