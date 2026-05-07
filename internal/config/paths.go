package config

import (
	"os"
	"path/filepath"
)

func UserConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "qsearch")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(home, ".config", "qsearch")
}

func UserProvidersPath() string {
	return filepath.Join(UserConfigDir(), "providers.toml")
}
