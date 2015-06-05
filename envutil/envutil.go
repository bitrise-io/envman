package envutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitrise-io/envman/pathutil"
	"gopkg.in/yaml.v2"
)

type EnvYMLStruct struct {
	Key   string `yml:"key"`
	Value string `yml:"value"`
}

type EnvListYMLStruct struct {
	Envlist []EnvYMLStruct `yml:"environment_variables"`
}

func ReadEnvListFromFile(path string) (EnvListYMLStruct, error) {
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		fmt.Println("Failed to check path, err: %s", err)
		return EnvListYMLStruct{}, err
	}
	if isExists == false {
		return EnvListYMLStruct{}, errors.New("No environemt variable list found")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return EnvListYMLStruct{}, err
	}

	var envlist EnvListYMLStruct
	err = yaml.Unmarshal(bytes, &envlist)
	if err != nil {
		return EnvListYMLStruct{}, err
	}

	return envlist, nil
}

func generateFormattedYMLForEnvList(envlist EnvListYMLStruct) ([]byte, error) {
	bytes, err := yaml.Marshal(envlist)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func WriteEnvListToFile(fpath string, envlist EnvListYMLStruct) error {
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
