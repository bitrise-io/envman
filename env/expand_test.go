package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitEnv(t *testing.T) {
	tests := []struct {
		name      string
		env       string
		wantKey   string
		wantValue string
	}{
		{
			name:      "simple case",
			env:       "A=B",
			wantKey:   "A",
			wantValue: "B",
		},
		{
			name:      "equals sign",
			env:       "A==B",
			wantKey:   "A",
			wantValue: "=B",
		},
		{
			name:      "",
			env:       "A=B=C=D",
			wantKey:   "A",
			wantValue: "B=C=D",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue := SplitEnv(tt.env)
			require.Equal(t, tt.wantKey, gotKey, "parseOSEnvs() gotKey")
			require.Equal(t, tt.wantValue, gotValue, "parseOSEnvs() gotvalue")
		})
	}
}
