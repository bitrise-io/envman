package cli

import (
	"os"
	"path"
	"testing"

	"github.com/bitrise-io/envman/envman"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	envsStr := `
envs:
- TEST_HOME1: $HOME
- TEST_HOME2: $TEST_HOME1/test
`
	environments, err := envman.ParseEnvsYML([]byte(envsStr))
	require.Equal(t, nil, err)

	envsJSONList, err := convertToEnsJSONModel(environments, false, false)
	require.Equal(t, nil, err)
	require.Equal(t, "$HOME", envsJSONList["TEST_HOME1"])
	require.Equal(t, "$TEST_HOME1/test", envsJSONList["TEST_HOME2"])

	testHome1 := os.Getenv("HOME")
	testHome2 := path.Join(testHome1, "test")
	envsJSONList, err = convertToEnsJSONModel(environments, true, false)
	require.Equal(t, nil, err)
	require.Equal(t, testHome1, envsJSONList["TEST_HOME1"])
	require.Equal(t, testHome2, envsJSONList["TEST_HOME2"])
}

func TestPrint_Sensitive(t *testing.T) {
	envsStr := `
envs:
- nonsensitivekey: testvalue
- sensitivekey: testsecret
  opts:
    is_sensitive: true
`
	environments, err := envman.ParseEnvsYML([]byte(envsStr))
	require.Equal(t, nil, err)

	// print everything
	envsJSONList, err := convertToEnsJSONModel(environments, false, false)
	require.Equal(t, nil, err)
	require.Equal(t, "testvalue", envsJSONList["nonsensitivekey"])
	require.Equal(t, "testsecret", envsJSONList["sensitivekey"])

	// print sensitive only
	envsJSONList, err = convertToEnsJSONModel(environments, true, true)
	require.Equal(t, nil, err)
	require.Equal(t, "", envsJSONList["nonsensitivekey"])
	require.Equal(t, "testsecret", envsJSONList["sensitivekey"])
}
