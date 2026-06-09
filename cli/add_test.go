package cli

import (
	"strconv"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/v2/envman"
	"github.com/bitrise-io/envman/v2/models"
	"github.com/stretchr/testify/require"
)

func TestEnvListSizeInBytes(t *testing.T) {
	str100Bytes := strings.Repeat("a", 100)
	require.Equal(t, 100, len([]byte(str100Bytes)))

	env := models.EnvironmentItemModel{
		"key": str100Bytes,
	}

	envList := []models.EnvironmentItemModel{env}
	size, err := envListSizeInBytes(envList)
	require.Equal(t, nil, err)
	require.Equal(t, 100, size)

	envList = []models.EnvironmentItemModel{env, env}
	size, err = envListSizeInBytes(envList)
	require.Equal(t, nil, err)
	require.Equal(t, 200, size)
}

func TestValidateEnv(t *testing.T) {
	defaultConfig, err := envman.GetConfigs()
	require.NoError(t, err)

	tests := []struct {
		name    string
		key     string
		value   string
		envList []models.EnvironmentItemModel
		wantErr error
	}{
		{
			name:  "Max allowed env var value",
			key:   "key",
			value: strings.Repeat("a", defaultConfig.EnvBytesLimitInKB*1024),
		},
		{
			name:    "Max allowed env var list",
			key:     "key",
			value:   strings.Repeat("a", defaultConfig.EnvListBytesLimitInKB/2*1024),
			envList: []models.EnvironmentItemModel{{"key": strings.Repeat("a", defaultConfig.EnvListBytesLimitInKB/2*1024)}},
		},
		{
			name:    "Too big env var value",
			key:     "key",
			value:   strings.Repeat("a", defaultConfig.EnvBytesLimitInKB*1024+1),
			wantErr: NewEnvVarValueTooLargeError("key", float64(defaultConfig.EnvBytesLimitInKB)+(1.0/1024.0), float64(defaultConfig.EnvBytesLimitInKB)),
		},
		{
			name:    "Too big env var list",
			key:     "key",
			value:   "a",
			envList: []models.EnvironmentItemModel{{"key": strings.Repeat("a", defaultConfig.EnvListBytesLimitInKB*1024)}},
			wantErr: NewEnvVarListTooLargeError(((float64)(defaultConfig.EnvListBytesLimitInKB))+(float64)(len("a"))/1024.0, float64(defaultConfig.EnvListBytesLimitInKB)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validValue, err := validateEnv(tt.key, tt.value, tt.envList)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
				require.Equal(t, "", validValue)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.value, validValue)
			}
		})
	}
}

func TestValidateEnvWithEnvListOverride(t *testing.T) {
	defaultConfig, err := envman.GetConfigs()
	require.NoError(t, err)

	// A value bigger than the default per-env byte limit.
	bigValue := strings.Repeat("a", defaultConfig.EnvBytesLimitInKB*1024+1)

	t.Run("rejected without an override", func(t *testing.T) {
		_, err := validateEnv("key", bigValue, nil)
		require.Error(t, err)
	})

	t.Run("accepted when the override comes from envman's env list", func(t *testing.T) {
		// The override lives in the env list envman works with (the envstore), not in the
		// process environment, mirroring how it is set during a build.
		envList := []models.EnvironmentItemModel{
			{envman.EnvBytesLimitInKBEnvKey: strconv.Itoa(defaultConfig.EnvBytesLimitInKB * 2)},
			{envman.EnvListBytesLimitInKBEnvKey: strconv.Itoa(defaultConfig.EnvListBytesLimitInKB * 2)},
		}

		validValue, err := validateEnv("key", bigValue, envList)
		require.NoError(t, err)
		require.Equal(t, bigValue, validValue)
	})

	t.Run("invalid override value in the env list returns an error", func(t *testing.T) {
		envList := []models.EnvironmentItemModel{
			{envman.EnvBytesLimitInKBEnvKey: "not-a-number"},
		}

		_, err := validateEnv("key", "a", envList)
		require.Error(t, err)
	})
}
