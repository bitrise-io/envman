package integration

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/envman"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func addCommand(key, value, envstore string) *command.Model {
	return command.New(binPath(), "-l", "debug", "-p", envstore, "add", "--key", key, "--value", value)
}

func addFileCommand(key, pth, envstore string) *command.Model {
	return command.New(binPath(), "-l", "debug", "-p", envstore, "add", "--key", key, "--valuefile", pth)
}

func addPipeCommand(key string, reader io.Reader, envstore string) *command.Model {
	return command.New(binPath(), "-l", "debug", "-p", envstore, "add", "--key", key).SetStdin(reader)
}

func runWithCustomEnvmanConfig(cfg envman.ConfigsModel, fn func()) error {
	configPath := filepath.Join(pathutil.UserHomeDir(), ".envman", "configs.json")
	if err := pathutil.EnsureDirExist(filepath.Dir(configPath)); err != nil {
		return err
	}

	exists, err := pathutil.IsPathExists(configPath)
	if err != nil {
		return err
	}

	var origData []byte
	if exists {
		origData, err = fileutil.ReadBytesFromFile(configPath)
		if err != nil {
			return err
		}
	}

	cfgData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := fileutil.WriteBytesToFile(configPath, cfgData); err != nil {
		return err
	}

	fn()

	if exists {
		return fileutil.WriteBytesToFile(configPath, origData)
	}
	return os.RemoveAll(configPath)
}

func TestAdd(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__envman__")
	require.NoError(t, err)

	envstore := filepath.Join(tmpDir, ".envstore")
	f, err := os.Create(envstore)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	t.Log("add flag value")
	{
		out, err := addCommand("KEY", "value", envstore).RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		cont, err := fileutil.ReadStringFromFile(envstore)
		require.NoError(t, err, out)
		require.Equal(t, "envs:\n- KEY: value\n", cont)
	}

	t.Log("add file flag value")
	{
		pth := filepath.Join(tmpDir, "file")
		require.NoError(t, fileutil.WriteStringToFile(pth, "some content"))

		out, err := addFileCommand("KEY", pth, envstore).RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		cont, err := fileutil.ReadStringFromFile(envstore)
		require.NoError(t, err, out)
		require.Equal(t, "envs:\n- KEY: some content\n", cont)
	}

	t.Log("add piped value")
	{
		out, err := addPipeCommand("KEY", strings.NewReader("some piped value"), envstore).RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		cont, err := fileutil.ReadStringFromFile(envstore)
		require.NoError(t, err, out)
		require.Equal(t, "envs:\n- KEY: some piped value\n", cont)
	}

	t.Log("add piped value over limit")
	{
		require.NoError(t, runWithCustomEnvmanConfig(envman.ConfigsModel{EnvBytesLimitInKB: 1, EnvListBytesLimitInKB: 2}, func() {
			out, err := addPipeCommand("KEY", strings.NewReader(strings.Repeat("0", 2*1024)), envstore).RunAndReturnTrimmedCombinedOutput()
			require.Error(t, err, out)
		}))
	}

	t.Log("add empty piped value")
	{
		out, err := addPipeCommand("KEY", strings.NewReader(""), envstore).RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		cont, err := fileutil.ReadStringFromFile(envstore)
		require.NoError(t, err, out)
		require.Equal(t, "envs:\n- KEY: \"\"\n", cont)
	}
}
