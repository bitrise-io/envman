package cli

import (
	"strconv"
	"strings"
	"testing"

	"github.com/bitrise-io/envman/v2/envman"
	"github.com/bitrise-io/envman/v2/models"
	"github.com/stretchr/testify/require"
)

func TestValidateEnvList(t *testing.T) {
	defaultConfig, err := envman.GetConfigs()
	require.NoError(t, err)

	tests := []struct {
		name    string
		envList []models.EnvironmentItemModel
		wantErr error
	}{
		{
			name:    "Max allowed env var value",
			envList: []models.EnvironmentItemModel{{"key": strings.Repeat("a", defaultConfig.EnvBytesLimitInKB*1024)}},
		},
		{
			name: "Max allowed env var list",
			envList: []models.EnvironmentItemModel{
				{"key1": strings.Repeat("a", defaultConfig.EnvListBytesLimitInKB/2*1024)},
				{"key2": strings.Repeat("a", defaultConfig.EnvListBytesLimitInKB/2*1024)},
			},
		},
		{
			name:    "Too big env var value",
			envList: []models.EnvironmentItemModel{{"key": strings.Repeat("a", defaultConfig.EnvBytesLimitInKB*1024+1)}},
			wantErr: NewEnvVarValueTooLargeError("key", float64(defaultConfig.EnvBytesLimitInKB)+(1.0/1024.0), float64(defaultConfig.EnvBytesLimitInKB)),
		},
		{
			name: "Too big env var list",
			envList: []models.EnvironmentItemModel{
				{"key1": strings.Repeat("a", defaultConfig.EnvListBytesLimitInKB*1024)},
				{"key2": "a"},
			},
			wantErr: NewEnvVarListTooLargeError(((float64)(defaultConfig.EnvListBytesLimitInKB))+(float64)(len("a"))/1024.0, float64(defaultConfig.EnvListBytesLimitInKB)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnvList(tt.envList)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateEnvListOverride(t *testing.T) {
	defaultConfig, err := envman.GetConfigs()
	require.NoError(t, err)

	// A value bigger than the default per-env byte limit.
	bigValue := strings.Repeat("a", defaultConfig.EnvBytesLimitInKB*1024+1)

	t.Run("rejected without an override", func(t *testing.T) {
		err := validateEnvList([]models.EnvironmentItemModel{{"key": bigValue}})
		require.Error(t, err)
	})

	t.Run("accepted when the override is anywhere in the list, even last", func(t *testing.T) {
		// The override keys are last in the list, after the value they need to permit. They are
		// pulled out before validation, so their position does not matter.
		envList := []models.EnvironmentItemModel{
			{"key": bigValue},
			{envman.EnvBytesLimitInKBEnvKey: strconv.Itoa(defaultConfig.EnvBytesLimitInKB * 2)},
			{envman.EnvListBytesLimitInKBEnvKey: strconv.Itoa(defaultConfig.EnvListBytesLimitInKB * 2)},
		}

		require.NoError(t, validateEnvList(envList))
	})

	t.Run("invalid override value returns an error", func(t *testing.T) {
		envList := []models.EnvironmentItemModel{
			{"key": "a"},
			{envman.EnvBytesLimitInKBEnvKey: "not-a-number"},
		}

		require.Error(t, validateEnvList(envList))
	})
}
