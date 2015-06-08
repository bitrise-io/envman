package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitrise-io/go-pathutil"
	"gopkg.in/yaml.v2"
)

type envYMLStruct struct {
	Key   string `yml:"key"`
	Value string `yml:"value"`
}

type envListYMLStruct struct {
	Envlist []envYMLStruct `yml:"environment_variables"`
}

func readEnvListFromFile(path string) (envListYMLStruct, error) {
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		fmt.Println("Failed to check path, err: %s", err)
		return envListYMLStruct{}, err
	}
	if isExists == false {
		return envListYMLStruct{}, errors.New("No environemt variable list found")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return envListYMLStruct{}, err
	}

	var envlist envListYMLStruct
	err = yaml.Unmarshal(bytes, &envlist)
	if err != nil {
		return envListYMLStruct{}, err
	}

	return envlist, nil
}

func generateFormattedYMLForEnvList(envlist envListYMLStruct) ([]byte, error) {
	bytes, err := yaml.Marshal(envlist)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func writeEnvListToFile(fpath string, envlist envListYMLStruct) error {
	if fpath == "" {
		return errors.New("No path provided")
	}

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonContBytes, err := generateFormattedYMLForEnvList(envlist)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonContBytes)
	if err != nil {
		return err
	}

	return nil
}
