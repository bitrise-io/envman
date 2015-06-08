package pathutil

import (
	"os"
	"runtime"
)

const envlistName string = "environment_variables.yml"

var envmanDir string = userHomeDir() + "/.envman/"

var EnvlistPath string = envmanDir + envlistName

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func CreateEnvmanDir() error {
	path := envmanDir
	exist, _ := IsPathExists(path)
	if exist {
		return nil
	}
	return createDir(path)
}

func createDir(path string) error {
	err := os.MkdirAll(path, 0755)
	return err
}
