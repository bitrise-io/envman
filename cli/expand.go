package cli

import (
	"fmt"
	"os"
	"strings"

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

func expandEnvironments(newEnvs []envmanModels.EnvironmentItemModel, initalEnvs map[string]envVarValue) (map[string]envVarValue, []declaredEnvVar, error) {
	envs := initalEnvs
	actionLog := make([]declaredEnvVar, len(newEnvs))

	// Expand enviroment variables, ordering of environments matters
	for i, env := range newEnvs {
		newEnv, err := declareEnvironmentVariable(env, envs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse new environment variable (%s): %s", env, err)
		}

		actionLog[i] = newEnv

		switch newEnv.action {
		case unsetEnv:
			delete(envs, newEnv.name)
		case skipEnv:
		case createEnv:
			envs[newEnv.name] = envVarValue{
				value:       newEnv.value,
				isSensitive: newEnv.isSensitive,
			}
		default:
			return nil, nil, fmt.Errorf("invalid case for environement declaration action: %#v", newEnv)
		}
	}

	return envs, actionLog, nil
}

func parseOSEnv(env string) (key string, value string) {
	const sep = "="
	split := strings.SplitAfterN(env, sep, 2)
	key = strings.TrimSuffix(split[0], sep)
	if len(split) > 1 {
		value = split[1]
	}
	return
}

func commandEnvs2(newEnvs []envmanModels.EnvironmentItemModel) ([]string, error) {
	initialOSEnvs := os.Environ()
	envs := make(map[string]envVarValue)

	for _, env := range initialOSEnvs {
		key, value := parseOSEnv(env)
		envs[key] = envVarValue{
			value:       value,
			isSensitive: false,
		}
	}

	_, actionLog, err := expandEnvironments(newEnvs, envs)
	if err != nil {
		return nil, err
	}

	for _, action := range actionLog {
		switch action.action {
		case createEnv:
			os.Setenv(action.name, action.value)
		case unsetEnv:
			os.Unsetenv(action.name)
		case skipEnv:
		default:
			return nil, fmt.Errorf("invalid case for environement declaration action: %#v", action)
		}
	}

	return os.Environ(), nil
}

func expandStepInputsForAnalytics(inputs, environments []envmanModels.EnvironmentItemModel) map[string]envVarValue {
	initialOSEnvs := os.Environ()
	initialEnvs := make(map[string]envVarValue)

	for _, env := range initialOSEnvs {
		key, value := parseOSEnv(env)
		initialEnvs[key] = envVarValue{
			value:       value,
			isSensitive: false,
		}
	}

	envs, _, err := expandEnvironments(append(environments, inputs...), initialEnvs)
	if err != nil {
		log.Warnf("%s", err)
	}

	// Filter inputs from enviroments
	expandedInputs := make(map[string]envVarValue)
	for _, input := range inputs {
		inputName, _, err := input.GetKeyValuePair()
		if err != nil {
			log.Warnf("failed to get new environment variable name and value: %s", err)
		}

		expandedInputs[inputName] = envs[inputName]
	}

	return expandedInputs
}
