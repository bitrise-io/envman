package integration

import (
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/bitrise-io/envman/env"
	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name         string
		envs, inputs []models.EnvironmentItemModel
		want         map[string]env.Variable
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

	tmpDir, err := pathutil.NormalizedOSTempDirPath("__envman__")
	require.NoError(t, err)

	envstore := filepath.Join(tmpDir, ".envstore")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnvmanInitAtPath(envstore)
			require.NoError(t, err, "EnvmanInitAtPath()")

			inEnvs := append(tt.envs, tt.inputs...)

			for _, envVar := range inEnvs {
				if err := envVar.FillMissingDefaults(); err != nil {
					require.NoError(t, err, "FillMissingDefaults()")
				}
			}

			ExportEnvironmentsList(envstore, inEnvs)
			require.NoError(t, err, "ExportEnvironmentsList()")

			output, err := EnvmanRun(envstore, tmpDir, []string{"env"}, time.Minute, nil)
			require.NoError(t, err, "EnvmanRun()")

			gotOut := parseEnvRawOut(output)

			t.Logf("Actual environment: %s", gotOut)

			// Add inital envrionments to the want map
			wantEnvs := make(map[string]string)
			for changedKey, changedValue := range tt.want {
				wantEnvs[changedKey] = changedValue.Value
			}

			for wantKey, wantValue := range wantEnvs {
				require.Equal(t, wantValue, gotOut[wantKey], "Set environments do not match.")
			}
		})
	}

}

func parseEnvRawOut(output string) map[string]string {
	r := regexp.MustCompile("(?m)^([a-zA-Z_][a-zA-Z0-9_]*)=(.*)$")
	matches := r.FindAllStringSubmatch(output, -1)

	envs := make(map[string]string)
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
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
A=`,
			want: map[string]string{
				"RBENV_SHELL": "zsh",
				"_":           "/usr/local/bin/go",
				"A":           "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseEnvRawOut(tt.output); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEnvRawOut() = %v, want %v", got, tt.want)
			}
		})
	}
}
