package config

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed defaults/providers.toml
var DefaultProvidersTOML string

func InitUserProvidersFile(force bool) (string, error) {
	path := UserProvidersPath()

	if _, err := os.Stat(path); err == nil && !force {
		return path, nil
	} else if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	if err := os.WriteFile(path, []byte(DefaultProvidersTOML), 0o644); err != nil {
		return "", err
	}

	return path, nil
}

func EnsureUserProvidersFile() (string, error) {
	return InitUserProvidersFile(false)
}
