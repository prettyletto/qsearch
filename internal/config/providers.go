package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/prettyletto/qsearch/internal/domain/provider"
	"github.com/prettyletto/qsearch/internal/providers/custom"
)

type ProvidersFile struct {
	Providers []ProviderConfig `toml:"providers"`
}

type ProviderConfig struct {
	Name      string   `toml:"name"`
	Aliases   []string `toml:"aliases"`
	URL       string   `toml:"url"`
	Icon      string   `toml:"icon"`
	TagBG     string   `toml:"tag_bg"`
	IconColor string   `toml:"icon_color"`
	TextColor string   `toml:"text_color"`
}

func LoadProviders(path string) ([]provider.Provider, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var file ProvidersFile
	if err := toml.Unmarshal(body, &file); err != nil {
		return nil, err
	}

	if len(file.Providers) == 0 {
		return nil, fmt.Errorf("no providers found in %s", path)
	}

	providers := make([]provider.Provider, 0, len(file.Providers))

	for _, cfg := range file.Providers {
		p, err := providerFromConfig(cfg)
		if err != nil {
			return nil, err
		}

		providers = append(providers, p)
	}
	return providers, nil
}

func providerFromConfig(cfg ProviderConfig) (provider.Provider, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("provider name is required")
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("provider %q url is required", cfg.Name)
	}

	if !strings.Contains(cfg.URL, "{{query}}") {
		return nil, fmt.Errorf("provider %q url must contain {{query}}", cfg.Name)
	}

	names := append([]string{cfg.Name}, cfg.Aliases...)

	meta := provider.Meta{
		Icon:      cfg.Icon,
		Name:      cfg.Name,
		TagBG:     cfg.TagBG,
		IconColor: cfg.IconColor,
		TextColor: cfg.TextColor,
	}

	return custom.New(names, cfg.URL, meta), nil
}
