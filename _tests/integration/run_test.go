package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/env"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__envman__")
	require.NoError(t, err)

	envstore := filepath.Join(tmpDir, ".envstore")

	for _, tt := range SharedTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := EnvmanInitAtPath(envstore)
			require.NoError(t, err, "EnvmanInitAtPath()")

			for _, envVar := range tt.Envs {
				if err := envVar.FillMissingDefaults(); err != nil {
					require.NoError(t, err, "FillMissingDefaults()")
				}
			}

			ExportEnvironmentsList(envstore, tt.Envs)
			require.NoError(t, err, "ExportEnvironmentsList()")

			output, err := EnvmanRun(envstore, tmpDir, []string{"env"})
			require.NoError(t, err, "EnvmanRun()")

			gotOut := parseEnvRawOut(output)

			// Want envs
			envsWant := make(map[string]string)
			for _, envVar := range os.Environ() {
				key, value := env.SplitEnv(envVar)
				envsWant[key] = value
			}

			for _, envCommand := range tt.Want {
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
