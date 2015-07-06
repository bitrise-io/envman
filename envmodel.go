package main

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
)

const (
	IS_EXPAND_KEY string = "is_expand"
	TRUE_KEY      string = "true"
	FALSE_KEY     string = "false"
)

/*
	This is the model of ENVIRONMENT in envman, for methods
*/
type envModel struct {
	Key      string
	Value    string
	IsExpand bool
}

/*
	This is the model of ENVIRONMENT in envman, for storing in file
*/
type envMap map[string]string

type envsYMLModel struct {
	Envs []envMap `yml:"environments"`
}

/*
	Convert envsYMLModel to envModel array
*/
func convertToEnvModelArray(envYML envsYMLModel) []envModel {
	var envModels []envModel

	for _, envMap := range envYML.Envs {
		envModel := convertToEnvModel(envMap)
		envModels = append(envModels, envModel)
	}

	return envModels
}

func convertToEnvModel(eMap envMap) envModel {
	var eModel envModel

	for key, value := range eMap {
		if key != IS_EXPAND_KEY {
			eModel.Key = key
			eModel.Value = value
		}
	}

	eModel.IsExpand = isExpand(eMap[IS_EXPAND_KEY])

	return eModel
}

func isExpand(s string) bool {
	if s == "" {
		return true
	} else {
		expand, err := strconv.ParseBool(s)
		if err != nil {
			log.Errorln("isExpand: Failed to parse input:", err)
			return true
		}
		return expand
	}
}

/*
	Convert envModel array to envsYMLModel
*/
func convertToEnvsYMLModel(eModels []envModel) envsYMLModel {
	var envYML envsYMLModel
	var envMaps []envMap

	for _, eModel := range eModels {
		eMap := convertToEnvMap(eModel)
		envMaps = append(envMaps, eMap)
	}

	envYML.Envs = envMaps
	return envYML
}

func convertToEnvMap(eModel envModel) envMap {
	eMap := make(envMap)

	if eModel.IsExpand == false {
		eMap[IS_EXPAND_KEY] = FALSE_KEY
	}

	eMap[eModel.Key] = eModel.Value

	return eMap
}
