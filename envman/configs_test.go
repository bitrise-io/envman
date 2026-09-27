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

func TestGetConfigsEnvOverride(t *testing.T) {
	fakeHomePth, err := pathutil.NormalizedOSTempDirPath("_FAKE_HOME")
	require.NoError(t, err)
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		require.NoError(t, os.Setenv("HOME", originalHome))
		require.NoError(t, os.RemoveAll(fakeHomePth))
	})
	require.NoError(t, os.Setenv("HOME", fakeHomePth))

	setEnvs := func(t *testing.T, envs map[string]string) {
		for _, key := range []string{EnvBytesLimitInKBEnvKey, EnvListBytesLimitInKBEnvKey} {
			require.NoError(t, os.Unsetenv(key))
		}
		for key, value := range envs {
			require.NoError(t, os.Setenv(key, value))
		}
		t.Cleanup(func() {
			for key := range envs {
				require.NoError(t, os.Unsetenv(key))
			}
		})
	}

	t.Run("env override applies over defaults when no config file exists", func(t *testing.T) {
		setEnvs(t, map[string]string{
			EnvBytesLimitInKBEnvKey:     "512",
			EnvListBytesLimitInKBEnvKey: "1024",
		})

		configs, err := GetConfigs()
		require.NoError(t, err)
		require.Equal(t, 512, configs.EnvBytesLimitInKB)
		require.Equal(t, 1024, configs.EnvListBytesLimitInKB)
	})

	t.Run("env override wins over configs.json", func(t *testing.T) {
		require.NoError(t, saveConfigs(ConfigsModel{EnvBytesLimitInKB: 123, EnvListBytesLimitInKB: 321}))
		t.Cleanup(func() { require.NoError(t, os.Remove(getEnvmanConfigsFilePath())) })

		setEnvs(t, map[string]string{EnvListBytesLimitInKBEnvKey: "1024"})

		configs, err := GetConfigs()
		require.NoError(t, err)
		require.Equal(t, 123, configs.EnvBytesLimitInKB, "unset override leaves the file value")
		require.Equal(t, 1024, configs.EnvListBytesLimitInKB, "set override replaces the file value")
	})

	t.Run("zero disables the limit", func(t *testing.T) {
		setEnvs(t, map[string]string{EnvListBytesLimitInKBEnvKey: "0"})

		configs, err := GetConfigs()
		require.NoError(t, err)
		require.Equal(t, 0, configs.EnvListBytesLimitInKB)
	})

	t.Run("malformed or negative override is ignored", func(t *testing.T) {
		for _, value := range []string{"not-a-number", "-1", "12kb", ""} {
			setEnvs(t, map[string]string{EnvListBytesLimitInKBEnvKey: value})

			configs, err := GetConfigs()
			require.NoError(t, err)
			require.Equal(t, defaultEnvListBytesLimitInKB, configs.EnvListBytesLimitInKB, "value %q should be ignored", value)
		}
	})
}

func TestSetConfigLimits(t *testing.T) {
	intPtr := func(i int) *int { return &i }

	setupFakeHome := func(t *testing.T) string {
		fakeHomePth, err := pathutil.NormalizedOSTempDirPath("_FAKE_HOME")
		require.NoError(t, err)
		originalHome := os.Getenv("HOME")
		t.Cleanup(func() {
			require.NoError(t, os.Setenv("HOME", originalHome))
			require.NoError(t, os.RemoveAll(fakeHomePth))
		})
		require.NoError(t, os.Setenv("HOME", fakeHomePth))
		// The env overrides must not leak in and mask the file value under test.
		for _, key := range []string{EnvBytesLimitInKBEnvKey, EnvListBytesLimitInKBEnvKey} {
			require.NoError(t, os.Unsetenv(key))
		}
		return fakeHomePth
	}

	t.Run("creates the file when none exists and restore removes it", func(t *testing.T) {
		setupFakeHome(t)
		configPth := getEnvmanConfigsFilePath()

		restore, err := SetConfigLimits(ConfigLimitOverrides{
			EnvBytesLimitInKB:     intPtr(512),
			EnvListBytesLimitInKB: intPtr(1024),
		})
		require.NoError(t, err)

		configs, err := GetConfigs()
		require.NoError(t, err)
		require.Equal(t, 512, configs.EnvBytesLimitInKB)
		require.Equal(t, 1024, configs.EnvListBytesLimitInKB)

		require.NoError(t, restore())
		exists, err := pathutil.IsPathExists(configPth)
		require.NoError(t, err)
		require.False(t, exists, "restore should remove a file that did not exist before")
	})

	t.Run("writes 0 explicitly to disable a limit", func(t *testing.T) {
		setupFakeHome(t)

		restore, err := SetConfigLimits(ConfigLimitOverrides{EnvListBytesLimitInKB: intPtr(0)})
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, restore()) })

		configs, err := GetConfigs()
		require.NoError(t, err)
		require.Equal(t, 0, configs.EnvListBytesLimitInKB, "0 must be persisted, not dropped to the default")
	})

	t.Run("merges with an existing file and preserves the untouched key", func(t *testing.T) {
		setupFakeHome(t)
		configPth := getEnvmanConfigsFilePath()

		require.NoError(t, saveConfigs(ConfigsModel{EnvBytesLimitInKB: 100, EnvListBytesLimitInKB: 200}))
		originalBytes, err := os.ReadFile(configPth)
		require.NoError(t, err)

		restore, err := SetConfigLimits(ConfigLimitOverrides{EnvListBytesLimitInKB: intPtr(1024)})
		require.NoError(t, err)

		configs, err := GetConfigs()
		require.NoError(t, err)
		require.Equal(t, 100, configs.EnvBytesLimitInKB, "untouched key should be preserved")
		require.Equal(t, 1024, configs.EnvListBytesLimitInKB)

		require.NoError(t, restore())
		restoredBytes, err := os.ReadFile(configPth)
		require.NoError(t, err)
		require.JSONEq(t, string(originalBytes), string(restoredBytes), "restore should return the exact original content")
	})

	t.Run("no overrides is a no-op and does not create a file", func(t *testing.T) {
		setupFakeHome(t)
		configPth := getEnvmanConfigsFilePath()

		restore, err := SetConfigLimits(ConfigLimitOverrides{})
		require.NoError(t, err)
		require.NoError(t, restore())

		exists, err := pathutil.IsPathExists(configPth)
		require.NoError(t, err)
		require.False(t, exists)
	})
}
