package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/_tests/integration"
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

func TestGetDeclarationsSideEffects(t *testing.T) {
	for _, test := range integration.SharedTestCases {
		t.Run(test.Name, func(t *testing.T) {
			// Arrange
			cleanEnvs := os.Environ()

			for _, envVar := range test.Envs {
				err := envVar.FillMissingDefaults()
				require.NoError(t, err, "FillMissingDefaults()")
			}
			// Act
			got, err := env.GetDeclarationsSideEffects(test.Envs, &env.DefaultEnvironmentSource{})
			require.NoError(t, err, "GetDeclarationsSideEffects()")

			err = restoreEnviron(cleanEnvs)
			require.NoError(t, err, "restoreEnviron()")

			// Assert
			require.NotNil(t, got)
			require.Equal(t, test.Want, got.CommandHistory)
		})
	}
}
