package cli

import "os"

// EnvController ...
type EnvController map[string]string

// Set ...
func (ec EnvController) Set(k, v string) error {
	ec[k] = v
	return nil
}

// Get ...
func (ec EnvController) Get(key string) string {
	return ec[key]
}

// Unset ...
func (ec EnvController) Unset(k string) error {
	delete(ec, k)
	return nil
}

// List ...
func (ec EnvController) List() (list []string) {
	for k, v := range ec {
		list = append(list, k+"="+v)
	}
	return
}

// Expand ...
func (ec EnvController) Expand(v string) string {
	return os.Expand(v, ec.Get)
}
