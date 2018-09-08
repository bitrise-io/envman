package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func binPath(t *testing.T) string {
	pth := os.Getenv("INTEGRATION_TEST_BINARY_PATH")
	require.NotEmpty(t, pth, "INTEGRATION_TEST_BINARY_PATH should not be empty")
	return pth
}
