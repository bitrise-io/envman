package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/env"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/stretchr/testify/require"
)

func TestExpandEnvsInString(t *testing.T) {
	t.Log("Expand env")
	{
		require.Equal(t, nil, os.Setenv("MY_ENV_KEY", "key"))

		inp := "${MY_ENV_KEY} of my home"
		expanded := expandEnvsInString(inp)

		key := os.Getenv("MY_ENV_KEY")
		require.NotEqual(t, "", expanded)
		require.Equal(t, key+" of my home", expanded)
	}
}

func TestCommandEnvs(t *testing.T) {
	t.Log("commandEnvs test")
	{
		env1 := models.EnvironmentItemModel{
			"test_key1": "test_value1",
		}
		require.Equal(t, nil, env1.FillMissingDefaults())

		env2 := models.EnvironmentItemModel{
			"test_key2": "test_value2",
		}
		require.Equal(t, nil, env2.FillMissingDefaults())

		envs := []models.EnvironmentItemModel{env1, env2}

		sessionEnvs, err := commandEnvs(envs)
		require.Equal(t, nil, err)

		env1Found := false
		env2Found := false
		for _, envString := range sessionEnvs {
			comp := strings.Split(envString, "=")
			key := comp[0]
			value := comp[1]

			envKey1, envValue1, err := env1.GetKeyValuePair()
			require.Equal(t, nil, err)

			envKey2, envValue2, err := env2.GetKeyValuePair()
			require.Equal(t, nil, err)

			if key == envKey1 && value == envValue1 {
				env1Found = true
			}
			if key == envKey2 && value == envValue2 {
				env2Found = true
			}
		}
		require.Equal(t, true, env1Found)
		require.Equal(t, true, env2Found)
	}

	// Test skip_if_empty
	t.Log("skip_if_empty=false && value=empty => should add")
	{
		env1 := models.EnvironmentItemModel{
			"test_key3": "",
		}
		require.Equal(t, nil, env1.FillMissingDefaults())

		env2 := models.EnvironmentItemModel{
			"test_key4": "test_value4",
		}
		require.Equal(t, nil, env2.FillMissingDefaults())

		envs := []models.EnvironmentItemModel{env1, env2}

		sessionEnvs, err := commandEnvs(envs)
		require.Equal(t, nil, err)

		env1Found := false
		env2Found := false
		for _, envString := range sessionEnvs {
			comp := strings.Split(envString, "=")
			key := comp[0]
			value := comp[1]

			envKey1, envValue1, err := env1.GetKeyValuePair()
			require.Equal(t, nil, err)

			envKey2, envValue2, err := env2.GetKeyValuePair()
			require.Equal(t, nil, err)

			if key == envKey1 && value == envValue1 {
				env1Found = true
			}
			if key == envKey2 && value == envValue2 {
				env2Found = true
			}
		}
		require.Equal(t, true, env1Found)
		require.Equal(t, true, env2Found)
	}

	t.Log("skip_if_empty=true && value=empty => should NOT add")
	{
		env1 := models.EnvironmentItemModel{
			"test_key5": "",
			"opts": models.EnvironmentItemOptionsModel{
				SkipIfEmpty: pointers.NewBoolPtr(true),
			},
		}
		require.Equal(t, nil, env1.FillMissingDefaults())

		env2 := models.EnvironmentItemModel{
			"test_key6": "test_value6",
		}
		require.Equal(t, nil, env2.FillMissingDefaults())

		envs := []models.EnvironmentItemModel{env1, env2}

		sessionEnvs, err := commandEnvs(envs)
		require.Equal(t, nil, err)

		env1Found := false
		env2Found := false
		for _, envString := range sessionEnvs {
			comp := strings.Split(envString, "=")
			key := comp[0]
			value := comp[1]

			envKey1, envValue1, err := env1.GetKeyValuePair()
			require.Equal(t, nil, err)

			envKey2, envValue2, err := env2.GetKeyValuePair()
			require.Equal(t, nil, err)

			if key == envKey1 && value == envValue1 {
				env1Found = true
			}
			if key == envKey2 && value == envValue2 {
				env2Found = true
			}
		}
		require.Equal(t, false, env1Found)
		require.Equal(t, true, env2Found)
	}

	t.Log("skip_if_empty=true && value=NOT_empty => should add")
	{
		env1 := models.EnvironmentItemModel{
			"test_key7": "test_value7",
			"opts": models.EnvironmentItemOptionsModel{
				SkipIfEmpty: pointers.NewBoolPtr(true),
			},
		}
		require.Equal(t, nil, env1.FillMissingDefaults())

		env2 := models.EnvironmentItemModel{
			"test_key8": "test_value8",
		}
		require.Equal(t, nil, env2.FillMissingDefaults())

		envs := []models.EnvironmentItemModel{env1, env2}

		sessionEnvs, err := commandEnvs(envs)
		require.Equal(t, nil, err)

		env1Found := false
		env2Found := false
		for _, envString := range sessionEnvs {
			comp := strings.Split(envString, "=")
			key := comp[0]
			value := comp[1]

			envKey1, envValue1, err := env1.GetKeyValuePair()
			require.Equal(t, nil, err)

			envKey2, envValue2, err := env2.GetKeyValuePair()
			require.Equal(t, nil, err)

			if key == envKey1 && value == envValue1 {
				env1Found = true
			}
			if key == envKey2 && value == envValue2 {
				env2Found = true
			}
		}
		require.Equal(t, true, env1Found)
		require.Equal(t, true, env2Found)
	}

	t.Log("expand envs test")
	{
		env1 := models.EnvironmentItemModel{
			"env1": "Hello",
		}
		require.Equal(t, nil, env1.FillMissingDefaults())

		env2 := models.EnvironmentItemModel{
			"env2": "${env1} world",
		}
		require.Equal(t, nil, env2.FillMissingDefaults())

		env3 := models.EnvironmentItemModel{
			"env3": "${env2} !",
		}
		require.Equal(t, nil, env3.FillMissingDefaults())

		envs := []models.EnvironmentItemModel{env1, env2, env3}

		sessionEnvs, err := commandEnvs(envs)
		require.Equal(t, nil, err)

		env3Found := false
		for _, envString := range sessionEnvs {
			comp := strings.Split(envString, "=")
			key := comp[0]
			value := comp[1]

			envKey3, _, err := env3.GetKeyValuePair()
			require.Equal(t, nil, err)

			if key == envKey3 {
				require.Equal(t, "Hello world !", value)
				env3Found = true
			}
		}
		require.Equal(t, true, env3Found)
	}

	t.Log("unset OS envs test")
	{
		// given
		key := "TEST_ENV"
		val := "test"
		if err := os.Setenv(key, val); err != nil {
			require.Equal(t, nil, err, "test setup: error seting env (%s=%s)", key, val)
		}
		env := models.EnvironmentItemModel{
			key: val,
			models.OptionsKey: models.EnvironmentItemOptionsModel{
				Unset: pointers.NewBoolPtr(true),
			},
		}
		require.Equal(t, nil, env.FillMissingDefaults())
		testEnvs := []models.EnvironmentItemModel{
			env,
		}

		// when
		envs, err := commandEnvs(testEnvs)
		envFmt := "%s=%s" // note: if this format mismatches elements of `envs`, test can be a false positive!
		unset := fmt.Sprintf(envFmt, key, val)

		// then
		require.Equal(t, nil, err)
		require.NotContains(t, envs, unset, "failed to unset env (%s)", key)

	}
}

func restoreEnviron(environ []string) error {
	currEnviron := os.Environ()
	for _, currEnv := range currEnviron {
		currEnvKey, _ := env.SplitEnv(currEnv)
		if err := os.Unsetenv(currEnvKey); err != nil {
			return err
		}
	}

	for _, envVar := range environ {
		key, value := env.SplitEnv(envVar)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set %s=%s: %s", key, value, err)
		}
	}

	return nil
}

// compare tests if each env var in environ is included in envs with the same value.
func compare(t *testing.T, environ []string, envs map[string]env.Variable) error {
	/*
		if len(environ) != len(envs) {
			return fmt.Errorf("compare() failed: elem num not equal (%d != %d)", len(environ), len(envs))
		}
	*/

	for _, envVar := range environ {
		key, value := env.SplitEnv(envVar)
		v, ok := envs[key]
		if ok != true {
			// return fmt.Errorf("compare() failed: %s not found", key)
		} else if v.Value != value {
			return fmt.Errorf("compare() failed: %s value (%s) not equals to: %s", key, value, v.Value)
		}
	}
	return nil
}

func TestExpandStepInputs(t *testing.T) {
	// Arrange
	tests := []struct {
		name   string
		envs   []models.EnvironmentItemModel
		inputs []models.EnvironmentItemModel
		want   map[string]env.Variable
	}{
		{
			name: "Env does not depend on input",
			envs: []models.EnvironmentItemModel{
				{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_major": "12", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_major":      {Value: "12"},
				"simulator_os_version": {Value: ""},
			},
		},
		{
			name: "Env does not depend on input (input order switched)",
			envs: []models.EnvironmentItemModel{
				{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_major": "12", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_major":      {Value: "12"},
				"simulator_os_version": {Value: ""},
			},
		},
		{
			name: "Secrets inputs are marked as sensitive",
			envs: []models.EnvironmentItemModel{},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"secret_input": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
				"secret_input":         {Value: "top secret", IsSensitive: true},
			},
		},
		{
			name: "Secrets environments are marked as sensitive",
			envs: []models.EnvironmentItemModel{
				{"secret_env": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_device": "iPhone $secret_env", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_device": {Value: "iPhone top secret", IsSensitive: true},
			},
		},
		{
			name: "Inputs referencing sensitive env are marked as sensitive",
			envs: []models.EnvironmentItemModel{
				{"date": "2020 $month"},
				{"month": "jun"},
				{"simulator_short_device": "($date) Ipad"},
				{"secret_simulator_device": "$simulator_short_device Pro", "opts": map[string]interface{}{"is_sensitive": true}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_device": "$secret_simulator_device"},
				{"fallback_simulator_device": "$simulator_device"},
			},
			want: map[string]env.Variable{
				"simulator_device":          {Value: "(2020 ) Ipad Pro", IsSensitive: true},
				"fallback_simulator_device": {Value: "(2020 ) Ipad Pro", IsSensitive: true},
			},
		},
		{
			name: "Not referencing other envs, missing options (sensive input).",
			envs: []models.EnvironmentItemModel{},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{}},
				{"simulator_device": "iPhone 8 Plus", "opts": map[string]interface{}{}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
				"simulator_device":     {Value: "iPhone 8 Plus"},
			},
		},
		{
			name: "Not referencing other envs, options specified.",
			envs: []models.EnvironmentItemModel{},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_device": "iPhone 8 Plus", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
				"simulator_device":     {Value: "iPhone 8 Plus"},
			},
		},
		{
			name: "Input references env var, is_expand is false.",
			envs: []models.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": false}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "$SIMULATOR_OS_VERSION"},
			},
		},
		{
			name: "Unset",
			envs: []models.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"unset": true}},
			},
			inputs: []models.EnvironmentItemModel{},
			want:   map[string]env.Variable{},
		},
		{
			name: "Skip if empty",
			envs: []models.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "", "opts": map[string]interface{}{"skip_if_empty": true}},
			},
			inputs: []models.EnvironmentItemModel{},
			want:   map[string]env.Variable{},
		},
		{
			name: "Env expansion, input contains env var.",
			envs: []models.EnvironmentItemModel{
				{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
			},
		},
		{
			name: "Env var expansion, input expansion",
			envs: []models.EnvironmentItemModel{
				{"SIMULATOR_OS_MAJOR_VERSION": "13", "opts": map[string]interface{}{"is_sensitive": false}},
				{"SIMULATOR_OS_MINOR_VERSION": "3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"SIMULATOR_OS_VERSION": "$SIMULATOR_OS_MAJOR_VERSION.$SIMULATOR_OS_MINOR_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
			},
		},
		{
			name: "Input expansion, input refers other input",
			envs: []models.EnvironmentItemModel{
				{"simulator_os_version": "12.1", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
				"simulator_device":     {Value: "iPhone 8 (13.3)"},
			},
		},
		{
			name: "Input expansion, input can not refer other input declared after it",
			envs: []models.EnvironmentItemModel{
				{"simulator_os_version": "12.1", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{"is_sensitive": false}},
				{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"simulator_os_version": {Value: "13.3"},
				"simulator_device":     {Value: "iPhone 8 (12.1)"},
			},
		},
		{
			name: "Input refers itself, env refers itself",
			envs: []models.EnvironmentItemModel{
				{"ENV_LOOP": "$ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"loop": "$loop", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env_loop": "$ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"loop":     {Value: ""},
				"env_loop": {Value: ""},
			},
		},
		{
			name: "Input refers itself, env refers itself; both have prefix included",
			envs: []models.EnvironmentItemModel{
				{"ENV_LOOP": "Env Something: $ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"loop": "Something: $loop", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env_loop": "$ENV_LOOP", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"loop":     {Value: "Something: "},
				"env_loop": {Value: "Env Something: "},
			},
		},
		{
			name: "Inputs refer inputs in a chain, with prefix included",
			envs: []models.EnvironmentItemModel{},
			inputs: []models.EnvironmentItemModel{
				{"similar2": "anything", "opts": map[string]interface{}{"is_sensitive": false}},
				{"similar": "$similar2", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env": "Something: $similar", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"similar2": {Value: "anything"},
				"similar":  {Value: "anything"},
				"env":      {Value: "Something: anything"},
			},
		},
		{
			name: "References in a loop are not expanded",
			envs: []models.EnvironmentItemModel{
				{"B": "$A", "opts": map[string]interface{}{"is_sensitive": false}},
				{"A": "$B", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			inputs: []models.EnvironmentItemModel{
				{"a": "$b", "opts": map[string]interface{}{"is_sensitive": false}},
				{"b": "$c", "opts": map[string]interface{}{"is_sensitive": false}},
				{"c": "$a", "opts": map[string]interface{}{"is_sensitive": false}},
				{"env": "$A", "opts": map[string]interface{}{"is_sensitive": false}},
			},
			want: map[string]env.Variable{
				"a":   {Value: ""},
				"b":   {Value: ""},
				"c":   {Value: ""},
				"env": {Value: ""},
			},
		},
	}

	for _, test := range tests {
		// t.Run(test.name, func(t *testing.T) {
		// 	cleanEnvs := os.Environ()
		// 	// Act
		// 	got := expandStepInputsForAnalytics(test.inputs, test.envs)
		// 	if err := restoreEnviron(cleanEnvs); err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	// envman expand
		// 	// environ, err := filterSecrets(append(test.envs, test.inputs...))
		// 	// require.NoError(t, err)
		// 	environ := append(test.envs, test.inputs...)
		// 	envmanEnvs, err := cli.commandEnvs(environ)
		// 	require.NoError(t, err)
		// 	if err := restoreEnviron(cleanEnvs); err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	// Assert
		// 	require.NotNil(t, got)
		// 	//require.NotNil(t, envmanEnvs)
		// 	if !reflect.DeepEqual(test.want, got) {
		// 		t.Fatalf("expandStepInputs() actual: %v expected: %v", got, test.want)
		// 	}
		// 	// Test if the test expectation is align with envman.CommandEnvs function's behaviour
		// 	if err := compare(t, envmanEnvs, test.want); err != nil {
		// 		t.Fatal(err)
		// 	}
		// 	// Test if expandStepInputsForAnalytics is align with envman.CommandEnvs function's behaviour
		// 	if err := compare(t, envmanEnvs, got); err != nil {
		// 		t.Fatal(err)
		// 	}
		// })

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
