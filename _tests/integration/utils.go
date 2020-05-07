package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
)

// EnvmanInitAtPath ...
func EnvmanInitAtPath(envstorePth string) error {
	const logLevel = "debug"
	args := []string{"--loglevel", logLevel, "--path", envstorePth, "init", "--clear"}
	return command.RunCommand(binPath(), args...)
}

// EnvmanAdd ...
func EnvmanAdd(envstorePth, key, value string, expand, skipIfEmpty bool) error {
	const logLevel = "debug"
	args := []string{"--loglevel", logLevel, "--path", envstorePth, "add", "--key", key, "--append"}
	if !expand {
		args = append(args, "--no-expand")
	}
	if skipIfEmpty {
		args = append(args, "--skip-if-empty")
	}

	envman := exec.Command(binPath(), args...)
	envman.Stdin = strings.NewReader(value)
	envman.Stdout = os.Stdout
	envman.Stderr = os.Stderr
	return envman.Run()
}

// EnvmanAdd ...
func EnvmanUnset(envstorePth, key, value string, expand, skipIfEmpty bool) error {
	const logLevel = "debug"
	args := []string{"--loglevel", logLevel, "--path", envstorePth, "unset", "--key", key /*"--append"*/}
	if !expand {
		args = append(args, "--no-expand")
	}
	if skipIfEmpty {
		args = append(args, "--skip-if-empty")
	}

	envman := exec.Command(binPath(), args...)
	envman.Stdin = strings.NewReader(value)
	envman.Stdout = os.Stdout
	envman.Stderr = os.Stderr
	return envman.Run()
}

// ExportEnvironmentsList ...
func ExportEnvironmentsList(envstorePth string, envsList []models.EnvironmentItemModel) error {
	for _, env := range envsList {
		key, value, err := env.GetKeyValuePair()
		if err != nil {
			return err
		}

		opts, err := env.GetOptions()
		if err != nil {
			return err
		}

		isExpand := models.DefaultIsExpand
		if opts.IsExpand != nil {
			isExpand = *opts.IsExpand
		}

		skipIfEmpty := models.DefaultSkipIfEmpty
		if opts.SkipIfEmpty != nil {
			skipIfEmpty = *opts.SkipIfEmpty
		}

		if opts.Unset != nil && *opts.Unset {
			if err := EnvmanUnset(envstorePth, key, value, isExpand, skipIfEmpty); err != nil {
				return err
			}
			return nil
		}

		if err := EnvmanAdd(envstorePth, key, value, isExpand, skipIfEmpty); err != nil {
			return err
		}
	}
	return nil
}

// EnvmanClear ...
func EnvmanClear(envstorePth string) error {
	const logLevel = "debug"
	args := []string{"--loglevel", logLevel, "--path", envstorePth, "clear"}
	out, err := command.New(binPath(), args...).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		errorMsg := err.Error()
		if errorutil.IsExitStatusError(err) && out != "" {
			errorMsg = out
		}
		return fmt.Errorf("failed to clear envstore (%s), error: %s", envstorePth, errorMsg)
	}
	return nil
}

// EnvmanRun runs a command through envman.
func EnvmanRun(envstorePth, workDir string, cmdArgs []string) (string, error) {
	const logLevel = "panic"
	args := []string{"--loglevel", logLevel, "--path", envstorePth, "run"}
	args = append(args, cmdArgs...)

	cmd := command.New(binPath(), args...).SetDir(workDir)

	return cmd.RunAndReturnTrimmedCombinedOutput()
}
