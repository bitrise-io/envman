package envman

import (
	"os"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestGetConfigs(t *testing.T) {
	// fake home, to save the configs into
	fakeHomePth, err := pathutil.NormalizedOSTempDirPath("_FAKE_HOME")
	t.Logf("fakeHomePth: %s", fakeHomePth)
	require.NoError(t, err)
	originalHome := os.Getenv("HOME")
	defer func() {
		require.NoError(t, os.Setenv("HOME", originalHome))
		require.NoError(t, os.RemoveAll(fakeHomePth))
	}()
	require.Equal(t, nil, os.Setenv("HOME", fakeHomePth))

	configPth := getEnvmanConfigsFilePath()
	t.Logf("configPth: %s", configPth)

	// --- TESTING

	baseConf, err := GetConfigs()
	t.Logf("baseConf: %#v", baseConf)
	require.NoError(t, err)
	require.Equal(t, defaultEnvBytesLimitInKB, baseConf.EnvBytesLimitInKB)
	require.Equal(t, defaultEnvListBytesLimitInKB, baseConf.EnvListBytesLimitInKB)

	// modify it
	baseConf.EnvBytesLimitInKB = 123
	baseConf.EnvListBytesLimitInKB = 321

	// save to file
	require.NoError(t, saveConfigs(baseConf))

	// read it back
	configs, err := GetConfigs()
	t.Logf("configs: %#v", configs)
	require.NoError(t, err)
	require.Equal(t, configs, baseConf)
	require.Equal(t, 123, configs.EnvBytesLimitInKB)
	require.Equal(t, 321, configs.EnvListBytesLimitInKB)

	// delete the tmp config file
	require.NoError(t, os.Remove(configPth))
}

func TestGetConfigsIgnoresProcessEnv(t *testing.T) {
	// The byte-limit overrides are read from envman's env list during validation, not from the
	// process environment, so GetConfigs must ignore these env vars.
	fakeHomePth, err := pathutil.NormalizedOSTempDirPath("_FAKE_HOME")
	require.NoError(t, err)
	originalHome := os.Getenv("HOME")
	defer func() {
		require.NoError(t, os.Setenv("HOME", originalHome))
		require.NoError(t, os.RemoveAll(fakeHomePth))
	}()
	require.NoError(t, os.Setenv("HOME", fakeHomePth))

	require.NoError(t, os.Setenv(EnvBytesLimitInKBEnvKey, "111"))
	require.NoError(t, os.Setenv(EnvListBytesLimitInKBEnvKey, "222"))
	defer func() {
		require.NoError(t, os.Unsetenv(EnvBytesLimitInKBEnvKey))
		require.NoError(t, os.Unsetenv(EnvListBytesLimitInKBEnvKey))
	}()

	configs, err := GetConfigs()
	require.NoError(t, err)
	require.Equal(t, defaultEnvBytesLimitInKB, configs.EnvBytesLimitInKB)
	require.Equal(t, defaultEnvListBytesLimitInKB, configs.EnvListBytesLimitInKB)
}

func TestLimitFromEnvValue(t *testing.T) {
	limit, err := LimitFromEnvValue("512", EnvBytesLimitInKBEnvKey)
	require.NoError(t, err)
	require.Equal(t, 512, limit)

	_, err = LimitFromEnvValue("not-a-number", EnvBytesLimitInKBEnvKey)
	require.Error(t, err)
}
