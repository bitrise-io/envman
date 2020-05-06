package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bitrise-io/envman/env"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

var SharedTestCases = []struct {
	name string
	envs []models.EnvironmentItemModel
	want []env.Command
}{
	{
		name: "empty env list",
		envs: []models.EnvironmentItemModel{},
		want: []env.Command{},
	},
	{
		name: "unset env",
		envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{"unset": true}},
		},
		want: []env.Command{
			{Action: env.UnsetAction, Variable: env.Variable{Key: "A"}},
		},
	},
	{
		name: "set env",
		envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
		},
	},
	{
		name: "set multiple envs",
		envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
			{"B": "C", "opts": map[string]interface{}{}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "B", Value: "C"}},
		},
	},
	{
		name: "set int env",
		envs: []models.EnvironmentItemModel{
			{"A": 12, "opts": map[string]interface{}{}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "12"}},
		},
	},
	{
		name: "skip env",
		envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
			{"S": "", "opts": map[string]interface{}{"skip_if_empty": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
			{Action: env.SkipAction, Variable: env.Variable{Key: "S"}},
		},
	},
	{
		name: "skip env, do not skip if not empty",
		envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
			{"S": "T", "opts": map[string]interface{}{"skip_if_empty": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "S", Value: "T"}},
		},
	},
	{
		name: "Env does only depend on envs declared before them",
		envs: []models.EnvironmentItemModel{
			{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_major": "12", "opts": map[string]interface{}{"is_expand": false}},
			{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_major", Value: "12"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: ""}},
		},
	},
	{
		name: "Env does only depend on envs declared before them (input order switched)",
		envs: []models.EnvironmentItemModel{
			{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_sensitive": false}},
			{"simulator_major": "12", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_major", Value: "12"}},
		},
	},
	{
		name: "Env does only depend on envs declared before them, envs in a loop",
		envs: []models.EnvironmentItemModel{
			{"A": "$C", "opts": map[string]interface{}{"is_expand": true}},
			{"B": "$A", "opts": map[string]interface{}{"is_expand": true}},
			{"C": "$B", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "B", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "C", Value: ""}},
		},
	},
	{
		name: "Do not expand env if is_expand is false",
		envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": false}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "$SIMULATOR_OS_VERSION"}},
		},
	},
	{
		name: "Expand env, self reference",
		envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_VERSION": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: ""}},
		},
	},
	{
		name: "Expand env, input contains env var",
		envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_expand": false}},
			{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
		},
	},
	{
		name: "Multi level env var expansion",
		envs: []models.EnvironmentItemModel{
			{"A": "1", "opts": map[string]interface{}{"is_expand": true}},
			{"B": "$A", "opts": map[string]interface{}{"is_expand": true}},
			{"C": "prefix $B", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "1"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "B", Value: "1"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "C", Value: "prefix 1"}},
		},
	},
	{
		name: "Multi level env var expansion 2",
		envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_MAJOR_VERSION": "13", "opts": map[string]interface{}{"is_expand": true}},
			{"SIMULATOR_OS_MINOR_VERSION": "3", "opts": map[string]interface{}{"is_expand": true}},
			{"SIMULATOR_OS_VERSION": "$SIMULATOR_OS_MAJOR_VERSION.$SIMULATOR_OS_MINOR_VERSION", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_MAJOR_VERSION", Value: "13"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_MINOR_VERSION", Value: "3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
		},
	},
	{
		name: "Env expansion (partial), step input can refers other input",
		envs: []models.EnvironmentItemModel{
			{"simulator_os_version": "13.3", "opts": map[string]interface{}{}},
			{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: "iPhone 8 (13.3)"}},
		},
	},
	{
		name: "Env expand, duplicate env declarations",
		envs: []models.EnvironmentItemModel{
			{"simulator_os_version": "12.1", "opts": map[string]interface{}{}},
			{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{"is_expand": "true"}},
			{"simulator_os_version": "13.3", "opts": map[string]interface{}{}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "12.1"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: "iPhone 8 (12.1)"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
		},
	},
	{
		name: "Secrets inputs are marked as sensitive",
		envs: []models.EnvironmentItemModel{
			{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			{"secret_input": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "secret_input", Value: "top secret", IsSensitive: true}},
		},
	},
	{
		name: "Input referencing secret env is marked as sensitive",
		envs: []models.EnvironmentItemModel{
			{"SECRET_ENV": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			{"simulator_device": "iPhone $SECRET_ENV", "opts": map[string]interface{}{"is_expand": true, "is_sensitive": false}},
		},
		want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SECRET_ENV", Value: "top secret", IsSensitive: true}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: "iPhone top secret", IsSensitive: true}},
		},
	},
}

func TestRun(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__envman__")
	require.NoError(t, err)

	envstore := filepath.Join(tmpDir, ".envstore")

	for _, tt := range SharedTestCases {
		t.Run(tt.name, func(t *testing.T) {
			err := EnvmanInitAtPath(envstore)
			require.NoError(t, err, "EnvmanInitAtPath()")

			for _, envVar := range tt.envs {
				if err := envVar.FillMissingDefaults(); err != nil {
					require.NoError(t, err, "FillMissingDefaults()")
				}
			}

			ExportEnvironmentsList(envstore, tt.envs)
			require.NoError(t, err, "ExportEnvironmentsList()")

			output, err := EnvmanRun(envstore, tmpDir, []string{"env"}, time.Minute, nil)
			require.NoError(t, err, "EnvmanRun()")

			gotOut := parseEnvRawOut(output)

			// Want envs
			envsWant := make(map[string]string)
			for _, envVar := range os.Environ() {
				key, value := env.SplitEnv(envVar)
				envsWant[key] = value
			}

			for _, envCommand := range tt.want {
				switch envCommand.Action {
				case env.SetAction:
					envsWant[envCommand.Variable.Key] = envCommand.Variable.Value
				case env.UnsetAction:
					delete(envsWant, envCommand.Variable.Key)
				case env.SkipAction:
				default:
					t.Fatalf("compare() failed, invalid action: %d", envCommand.Action)
				}
			}

			require.Equal(t, envsWant, gotOut)
		})
	}

}

func parseEnvRawOut(output string) map[string]string {
	r := regexp.MustCompile("^([a-zA-Z_][a-zA-Z0-9_]*)=(.*)$")

	lines := strings.Split(output, "\n")

	envs := make(map[string]string)
	lastKey := ""
	for _, line := range lines {
		match := r.FindStringSubmatch(line)

		fmt.Printf("%s %s \n", line, match)

		if match == nil {
			if lastKey != "" {
				envs[lastKey] += "\n" + line
			}
			continue
		}

		if len(match) != 3 {
			continue
		}

		if match[1] == "current_envman" {
			if lastKey != "" {
				envs[lastKey] += "\n" + line
			}
			continue
		}

		lastKey = match[1]
		envs[match[1]] = match[2]
	}

	return envs
}

func Test_parseEnvRawOut(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   map[string]string
	}{
		{
			output: `RBENV_SHELL=zsh
_=/usr/local/bin/go
#!/bin/env bash
echo "ff"
A=`,
			want: map[string]string{
				"RBENV_SHELL": "zsh",
				"_": `/usr/local/bin/go
#!/bin/env bash
echo "ff"`,
				"A": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseEnvRawOut(tt.output)
			require.Equal(t, got, tt.want)
		})
	}
}
