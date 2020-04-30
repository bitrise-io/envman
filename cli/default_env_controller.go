package cli

import "os"

// DefaultEnvController ...
type DefaultEnvController struct{}

// Set ...
func (DefaultEnvController) Set(k, v string) error {
	return os.Setenv(k, v)
}

// Unset ...
func (DefaultEnvController) Unset(k string) error {
	return os.Unsetenv(k)
}

// List ...
func (DefaultEnvController) List() []string {
	return os.Environ()
}

// Expand ...
func (DefaultEnvController) Expand(v string) string {
	return os.ExpandEnv(v)
}
