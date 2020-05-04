package cli

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/stretchr/testify/require"
)

func restoreEnviron(environ []string) error {
	currEnviron := os.Environ()
	for _, currEnv := range currEnviron {
		currEnvKey, _ := parseOSEnv(currEnv)
		if err := os.Unsetenv(currEnvKey); err != nil {
			return err
		}
	}

	for _, env := range environ {
		key, value := parseOSEnv(env)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set %s=%s: %s", key, value, err)
		}
	}

	return nil
}

func TestExpandStepInputs(t *testing.T) {
	// Arrange
	tests := []struct {
		name   string
		envs   []envmanModels.EnvironmentItemModel
		inputs []envmanModels.EnvironmentItemModel
		want   map[string]envVarValue
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
			want: map[string]envVarValue{
				"simulator_major":      {value: "12"},
				"simulator_os_version": {value: ""},
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
			want: map[string]envVarValue{
				"simulator_major":      {value: "12"},
				"simulator_os_version": {value: ""},
			},
		},
		{
			name: "Secrets inputs are marked as sensitive",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"secret_input": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
				"secret_input":         {value: "top secret", isSensitive: true},
			},
		},
		{
			name: "Secrets environments are marked as sensitive",
			envs: []envmanModels.EnvironmentItemModel{
				{"secret_env": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_device": "iPhone $secret_env", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]envVarValue{
				"simulator_device": {value: "iPhone top secret", isSensitive: true},
			},
		},
		{
			name: "Inputs referencing sensitive env are marked as sensitive",
			envs: []envmanModels.EnvironmentItemModel{
				{"date": "2020 $month"},
				{"month": "jun"},
				{"simulator_short_device": "($date) Ipad"},
				{"secret_simulator_device": "$simulator_short_device Pro", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_device": "$secret_simulator_device"},
				{"fallback_simulator_device": "$simulator_device"},
			},
			want: map[string]envVarValue{
				"simulator_device":          {value: "(2020 ) Ipad Pro", isSensitive: true},
				"fallback_simulator_device": {value: "(2020 ) Ipad Pro", isSensitive: true},
			},
		},
		{
			name: "Not referencing other envs, missing options (sensive input).",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{}},
				{"simulator_device": "iPhone 8 Plus", "opts": map[string]interface{}{}},
			},
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
				"simulator_device":     {value: "iPhone 8 Plus"},
			},
		},
		{
			name: "Not referencing other envs, options specified.",
			envs: []envmanModels.EnvironmentItemModel{},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_device": "iPhone 8 Plus", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
				"simulator_device":     {value: "iPhone 8 Plus"},
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
			want: map[string]envVarValue{
				"simulator_os_version": {value: "$SIMULATOR_OS_VERSION"},
			},
		},
		{
			name: "Unset",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"unset": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{},
			want:   map[string]envVarValue{},
		},
		{
			name: "Skip if empty",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "", "opts": map[string]interface{}{"skip_if_empty": true}},
			},
			inputs: []envmanModels.EnvironmentItemModel{},
			want:   map[string]envVarValue{},
		},
		{
			name: "Env expansion, input contains env var.",
			envs: []envmanModels.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []envmanModels.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
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
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
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
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
				"simulator_device":     {value: "iPhone 8 (13.3)"},
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
			want: map[string]envVarValue{
				"simulator_os_version": {value: "13.3"},
				"simulator_device":     {value: "iPhone 8 (12.1)"},
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
			want: map[string]envVarValue{
				"loop":     {value: ""},
				"env_loop": {value: ""},
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
			want: map[string]envVarValue{
				"loop":     {value: "Something: "},
				"env_loop": {value: "Env Something: "},
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
			want: map[string]envVarValue{
				"similar2": {value: "anything"},
				"similar":  {value: "anything"},
				"env":      {value: "Something: anything"},
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
			want: map[string]envVarValue{
				"a":   {value: ""},
				"b":   {value: ""},
				"c":   {value: ""},
				"env": {value: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanEnvs := os.Environ()
			// Act
			got := expandStepInputsForAnalytics(test.inputs, test.envs)
			if err := restoreEnviron(cleanEnvs); err != nil {
				t.Fatal(err)
			}

			// envman expand
			// environ, err := filterSecrets(append(test.envs, test.inputs...))
			// require.NoError(t, err)
			environ := append(test.envs, test.inputs...)
			envmanEnvs, err := commandEnvs(environ)
			require.NoError(t, err)
			if err := restoreEnviron(cleanEnvs); err != nil {
				t.Fatal(err)
			}

			// Assert
			require.NotNil(t, got)
			//require.NotNil(t, envmanEnvs)
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

		t.Run("Compare commandEnvs and commandEnvs2"+test.name, func(t *testing.T) {
			cleanEnvs := os.Environ()
			environ := append(test.envs, test.inputs...)

			for _, newEnv := range environ {
				if err := newEnv.FillMissingDefaults(); err != nil {
					t.Fatalf("failed to fill missing defaults: %s", err)
				}
			}

			got1, err := commandEnvs(environ)
			require.NoError(t, err)
			require.NotNil(t, got1)
			if err := restoreEnviron(cleanEnvs); err != nil {
				t.Fatal(err)
			}

			got2, err := commandEnvs2(environ)
			require.NoError(t, err)
			require.NotNil(t, got2)
			if err := restoreEnviron(cleanEnvs); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got1, got2) {
				t.Fatalf("commandEnvs2() actual: %#v, expecteed: %#v", got2, got1)
			}

			if compare(t, got2, test.want); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// compare tests if each env var in environ is included in envs with the same value.
func compare(t *testing.T, environ []string, envs map[string]envVarValue) error {
	/*
		if len(environ) != len(envs) {
			return fmt.Errorf("compare() failed: elem num not equal (%d != %d)", len(environ), len(envs))
		}
	*/

	for _, envVar := range environ {
		key, value := parseOSEnv(envVar)
		v, ok := envs[key]
		if ok != true {
			// return fmt.Errorf("compare() failed: %s not found", key)
		} else if v.value != value {
			return fmt.Errorf("compare() failed: %s value (%s) not equals to: %s", key, value, v.value)
		}
	}
	return nil
}
