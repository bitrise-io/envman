package main

/*
	This is the model of ENVIRONMENT in envman, for methods
*/
type envModel struct {
	Key   string
	Value string
}

/*
	This is the model of ENVIRONMENT in envman, for storing in file
*/
type envMap map[string]string

type envsYMLModel struct {
	Envs envMap `yml:"environments"`
}

/*
	Convert envsYMLModel to envModel array
*/
func convertToEnvModelArray(envYML envsYMLModel) []envModel {
	var envModels []envModel

	for key, value := range envYML.Envs {
		eModel := envModel{key, value}
		envModels = append(envModels, eModel)
	}

	return envModels
}

/*
	Convert envModel array to envsYMLModel
*/
func convertToEnvsYMLModel(eModels []envModel) envsYMLModel {
	var envYML envsYMLModel
	envYML.Envs = make(envMap)

	for _, eModel := range eModels {
		envYML.Envs[eModel.Key] = eModel.Value
	}

	return envYML
}
