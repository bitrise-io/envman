package cli

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	envmanModels "github.com/bitrise-io/envman/models"
)

type declaredEnvVarAction int

const (
	invalidAction declaredEnvVarAction = iota + 1
	unsetEnv
	skipEnv
	createEnv
)

type declaredEnvVar struct {
	action      declaredEnvVarAction
	name        string
	value       string
	isSensitive bool
}

type envVarValue struct {
	value       string
	isSensitive bool
}

func declareEnvironmentVariable(env envmanModels.EnvironmentItemModel, initalEnvs map[string]envVarValue) (declaredEnvVar, error) {
	if err := env.FillMissingDefaults(); err != nil {
		return declaredEnvVar{}, fmt.Errorf("failed to fill missing defaults: %s", err)
	}

	envName, envValue, err := env.GetKeyValuePair()
	if err != nil {
		return declaredEnvVar{}, fmt.Errorf("failed to get new environment variable name and value: %s", err)
	}

	options, err := env.GetOptions()
	if err != nil {
		return declaredEnvVar{}, fmt.Errorf("failed to get new environment options: %s", err)
	}

	if options.Unset != nil && *options.Unset {
		return declaredEnvVar{
			action: unsetEnv,
			name:   envName,
		}, nil
	}

	if options.SkipIfEmpty != nil && *options.SkipIfEmpty && envValue == "" {
		return declaredEnvVar{
			action: skipEnv,
			name:   envName,
		}, nil
	}

	mappingFuncFactory := func(envs map[string]envVarValue, isDeclaredEnvSensitive *bool) func(string) string {
		return func(key string) string {
			if _, ok := envs[key]; !ok {
				return ""
			}

			*isDeclaredEnvSensitive = *isDeclaredEnvSensitive || envs[key].isSensitive
			return envs[key].value
		}
	}

	isDeclaredEnvSensitive := *options.IsSensitive
	if options.IsExpand != nil && *options.IsExpand {
		envValue = os.Expand(envValue, mappingFuncFactory(initalEnvs, &isDeclaredEnvSensitive))
	}

	return declaredEnvVar{
		action:      createEnv,
		name:        envName,
		value:       envValue,
		isSensitive: isDeclaredEnvSensitive,
	}, nil
}

func expandStepInputsForAnalytics(inputs, environments []envmanModels.EnvironmentItemModel) map[string]envVarValue {
	envs := make(map[string]envVarValue)

	// Expand enviroment variables, ordering of environments matters
	for _, env := range environments {
		newEnv, err := declareEnvironmentVariable(env, envs)
		if err != nil {
			log.Warnf("Failed to handle new env variable (%s), skipping: %s", env, err)
			continue
		}

		switch newEnv.action {
		case unsetEnv:
			delete(envs, newEnv.name)
		case skipEnv:
			continue
		case createEnv:
			envs[newEnv.name] = envVarValue{
				value:       newEnv.value,
				isSensitive: newEnv.isSensitive,
			}
		}
	}

	expandedInputs := make(map[string]envVarValue)
	// Retrieve all non-sensitive input values and expand them, order of inputs matters
	for _, input := range inputs {
		newEnv, err := declareEnvironmentVariable(input, envs)
		if err != nil {
			log.Warnf("Failed to handle new input env variable (%s), skipping: %s", input, err)
			continue
		}

		switch newEnv.action {
		case unsetEnv:
			delete(envs, newEnv.name)
		case skipEnv:
			continue
		case createEnv:
			expandedInputs[newEnv.name] = envVarValue{} // Save input names, so we can filter from envs later
			envs[newEnv.name] = envVarValue{
				value:       newEnv.value,
				isSensitive: newEnv.isSensitive,
			}
		}
	}

	// Filter inputs from enviroments
	for inputName := range expandedInputs {
		expandedInputs[inputName] = envs[inputName]
	}

	return expandedInputs
}
