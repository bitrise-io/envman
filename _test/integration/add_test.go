package integration

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__envman__")
	require.NoError(t, err)

	envstorePth := filepath.Join(tmpDir, "envstore.yml")
	_, err = os.Create(envstorePth)
	require.NoError(t, err)

	t.Log("add using flag value")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST1", "--value", "test1")
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n", envstore)
	}

	t.Log("add using piped string")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST2").SetStdin(strings.NewReader("test2"))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n- TEST2: test2\n", envstore)
	}

	t.Log("add using flag and piped string - piped string has priority")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST2", "--value", "hello").SetStdin(strings.NewReader("test2"))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n- TEST2: test2\n", envstore)
	}

	t.Log("add using flag and piped empty string - empty string does not overrides the flag value")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST2", "--value", "hello").SetStdin(strings.NewReader(""))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n- TEST2: hello\n", envstore)
	}

	t.Log("add using piped empty string")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST3").SetStdin(strings.NewReader(""))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n- TEST2: hello\n- TEST3: \"\"\n", envstore)
	}

	t.Log("add using piped nil bytes")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST4").SetStdin(bytes.NewReader([]byte(nil)))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n- TEST2: hello\n- TEST3: \"\"\n- TEST4: \"\"\n", envstore)
	}

	t.Log("add using empty pipe")
	{
		cmd := command.New(binPath(t), "-l", "debug", "--path", envstorePth, "add", "--key", "TEST5")
		execCmd := cmd.GetCmd()
		_, err := execCmd.StdinPipe()
		require.NoError(t, err)

		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		envstore, err := fileutil.ReadStringFromFile(envstorePth)
		require.NoError(t, err)
		require.Equal(t, "envs:\n- TEST1: test1\n- TEST2: hello\n- TEST3: \"\"\n- TEST4: \"\"\n- TEST5: \"\"\n", envstore)
	}
}
