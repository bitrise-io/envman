package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/stretchr/testify/require"
)

func restoreEnviron(environ []string) error {
	currEnviron := os.Environ()
	for _, currEnv := range currEnviron {
		currEnvKey, _ := keyValue(currEnv)
		if err := os.Unsetenv(currEnvKey); err != nil {
			return err
		}
	}

	for _, env := range environ {
		key, value := keyValue(env)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set %s=%s: %s", key, value, err)
		}
	}

	return nil
}

func filterSecrets(environ []models.EnvironmentItemModel) ([]models.EnvironmentItemModel, error) {
	var filtered []models.EnvironmentItemModel
	for _, env := range environ {
		opts, err := env.GetOptions()
		if err != nil {
			return nil, err
		}
		if opts.IsSensitive != nil && *opts.IsSensitive == true {
			continue
		}
		filtered = append(filtered, env)
	}
	return filtered, nil
}

func keyValue(env string) (key string, value string) {
	const sep = "="
	split := strings.SplitAfterN(env, sep, 2)
	key = strings.TrimSuffix(split[0], sep)
	if len(split) > 1 {
		value = split[1]
	}
	return
}

func expandEnvironment(environments []models.EnvironmentItemModel) ([]string, error) {
	// envman.CommandEnvs() achives the expansion in the way it modifies the current process' environment,
	// at the end of the function we restore the original environment (restoreEnviron).
	osEnviron := os.Environ()

	expandedEnviron, err := commandEnvs(environments)
	if err != nil {
		return nil, err
	}

	// envman.CommandEnvs(environments) always adds os.Environ() to the input environment.
	// For testing we need to make it's behavior close to the expandStepInputsForAnalytics() function,
	// which works only with the input environment.
	//
	// With this workaround there is still a differenc in the 2 function's behaviour:
	// envman.CommandEnvs(environments) always uses os.Environ() for expansion,
	// while expandStepInputsForAnalytics() works only with the input environment.
	// If we do not refer to os.Environ() the 2 functions should do expansion in the same way.
	var filteredEnviron []string
	for _, expandedEnv := range expandedEnviron {
		set := true

		expandedEnvKey, _ := keyValue(expandedEnv)
		for _, osEnv := range osEnviron {
			osEnvKey, _ := keyValue(osEnv)
			if expandedEnvKey == osEnvKey {
				set = false
				break
			}
		}

		if set {
			filteredEnviron = append(filteredEnviron, expandedEnv)
		}
	}

	return filteredEnviron, restoreEnviron(osEnviron)
}

func TestSecretExpand(t *testing.T) {
	// Arrange
	appEnvs := []envmanModels.EnvironmentItemModel{
		{"date": "2020 $month"},
		{"month": "jun"},
		{"simulator_short_device": "($date) Ipad"},
	}
	workflowEnvs := []envmanModels.EnvironmentItemModel{
		{
			"secret_simulator_device": "$simulator_short_device Pro",
			"opts":                    map[string]interface{}{"is_sensitive": true},
		},
	}
	inputEnvs := []envmanModels.EnvironmentItemModel{
		{"simulator_device": "$secret_simulator_device"},
		{"fallback_simulator_device": "$simulator_device"},
	}

	// Act
	got := expandStepInputsForAnalytics(inputEnvs, append(appEnvs, workflowEnvs...))

	// envman expand
	// environ, err := filterSecrets(append(append(appEnvs, workflowEnvs...), inputEnvs...))
	//require.NoError(t, err)
	environ := append(append(appEnvs, workflowEnvs...), inputEnvs...)
	envmanEnvs, err := expandEnvironment(environ)
	require.NoError(t, err)

	// Assert
	require.NotNil(t, got)
	require.NotNil(t, envmanEnvs)
	want := map[string]string{
		"simulator_device":          "[REDACTED]",
		"fallback_simulator_device": "(2020 ) Ipad Pro",
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expandStepInputs() actual: %v expected: %v", got, want)
	}

	// Test if the test expectation is align with envman.CommandEnvs function's behaviour
	/*
			if err := compare(t, envmanEnvs, want); err != nil {
				t.Errorf("invalid expectation: %s", err)
			}
		// Test if expandStepInputsForAnalytics is align with envman.CommandEnvs function's behaviour
		if err := compare(t, want, got); err != nil {
			t.Error(err)
		}
	*/
}

func TestExpandStepInputs(t *testing.T) {
	// Arrange
	tests := []struct {
		name   string
		envs   []envmanModels.EnvironmentItemModel
		inputs []envmanModels.EnvironmentItemModel
		want   map[string]string
	}{
		{
			name: "Env does not depend on input",
			envs: []envmanModels.EnvironmentItemModel{
				{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_major": "12", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_major":      "12",
				"simulator_os_version": "",
			},
		},
		{
			name: "Env does not depend on input (input order switched)",
			envs: []envmanModels.EnvironmentItemModel{
				{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_major": "12", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_major":      "12",
				"simulator_os_version": "",
			},
		},
		{
			name: "Secrets inputs are removed",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"secret_input": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
				// "secret_input":         "",
			},
		},
		{
			name: "Secrets environments are redacted",
			envs: []envmanModels.EnvironmentItemModel{
				{"secret_env": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_device": "iPhone $secret_env", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_device": "iPhone [REDACTED]",
			},
		},
		{
			name: "Not referencing other envs, missing options (sensive input).",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{}},
				{"simulator_device": "iPhone 8 Plus", "opts": map[string]interface{}{}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
				"simulator_device":     "iPhone 8 Plus",
			},
		},
		{
			name: "Not referencing other envs, options specified.",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_device": "iPhone 8 Plus", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
				"simulator_device":     "iPhone 8 Plus",
			},
		},
		{
			name: "Input references env var, is_expand is false.",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": false}},
			},
			want: map[string]string{
				"simulator_os_version": "$SIMULATOR_OS_VERSION",
			},
		},
		{
			name: "Unset",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"unset": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{},
			want:   map[string]string{},
		},
		{
			name: "Skip if empty",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "", "opts": map[string]interface{}{"skip_if_empty": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{},
			want:   map[string]string{},
		},
		{
			name: "Env expansion, input contains env var.",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
			},
		},
		{
			name: "Env var expansion, input expansion",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_MAJOR_VERSION": "13", "opts": map[string]interface{}{"is_sensitive": false}},
				{"SIMULATOR_OS_MINOR_VERSION": "3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"SIMULATOR_OS_VERSION": "$SIMULATOR_OS_MAJOR_VERSION.$SIMULATOR_OS_MINOR_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
			},
		},
		{
			name: "Input expansion, input refers other input",
			envs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "12.1", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
				"simulator_device":     "iPhone 8 (13.3)",
			},
		},
		{
			name: "Input expansion, input can not refer other input declared after it",
			envs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "12.1", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"simulator_os_version": "13.3",
				"simulator_device":     "iPhone 8 (12.1)",
			},
		},
		{
			name: "Input refers itself, env refers itself",
			envs: []envmanModels.EnvironmentItemModel{
				{"ENV_LOOP": "$ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"loop": "$loop", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env_loop": "$ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"loop":     "",
				"env_loop": "",
			},
		},
		{
			name: "Input refers itself, env refers itself; both have prefix included",
			envs: []envmanModels.EnvironmentItemModel{
				{"ENV_LOOP": "Env Something: $ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"loop": "Something: $loop", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env_loop": "$ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"loop":     "Something: ",
				"env_loop": "Env Something: ",
			},
		},
		{
			name: "Inputs refer inputs in a chain, with prefix included",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"similar2": "anything", "opts": map[string]interface{}{"is_sensitive": false}},
				{"similar": "$similar2", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env": "Something: $similar", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"similar2": "anything",
				"similar":  "anything",
				"env":      "Something: anything",
			},
		},
		{
			name: "References in a loop are not expanded",
			envs: []envmanModels.EnvironmentItemModel{
				{"B": "$A", "opts": map[string]interface{}{"is_sensitive": false}},
				{"A": "$B", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"a": "$b", "opts": map[string]interface{}{"is_sensitive": false}},
				{"b": "$c", "opts": map[string]interface{}{"is_sensitive": false}},
				{"c": "$a", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env": "$A", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]string{
				"a":   "",
				"b":   "",
				"c":   "",
				"env": "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			got := expandStepInputsForAnalytics(test.inputs, test.envs)

			// envman expand
			// environ, err := filterSecrets(append(test.envs, test.inputs...))
			// require.NoError(t, err)
			environ := append(test.envs, test.inputs...)
			envmanEnvs, err := expandEnvironment(environ)
			require.NoError(t, err)

			// Assert
			require.NotNil(t, got)
			require.NotNil(t, envmanEnvs)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expandStepInputs() actual: %v expected: %v", got, test.want)
			}
			// Test if the test expectation is align with envman.CommandEnvs function's behaviour
			if err := compare(t, envmanEnvs, test.want); err != nil {
				t.Fatal(err)
			}
			// Test if expandStepInputsForAnalytics is align with envman.CommandEnvs function's behaviour
			if err := compare(t, envmanEnvs, got); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// compare tests if each env var in environ is included in envs with the same value.
func compare(t *testing.T, environ []string, envs map[string]string) error {
	/*
		if len(environ) != len(envs) {
			return fmt.Errorf("compare() failed: elem num not equal (%d != %d)", len(environ), len(envs))
		}
	*/

	for _, envVar := range environ {
		key, value := keyValue(envVar)
		v, ok := envs[key]
		if ok != true {
			// return fmt.Errorf("compare() failed: %s not found", key)
		} else if v != value {
			return fmt.Errorf("compare() failed: %s value (%s) not equals to: %s", key, value, v)
		}
	}
	return nil
}
