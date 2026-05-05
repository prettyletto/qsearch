package dispatch

import (
	"fmt"
	"qsearch/internal/domain/provider"
	"strings"
)

type SearchRunner interface {
	Run(p provider.Provider, args []string) error
}

type Dispatcher struct {
	providers map[string]provider.Provider
	search    SearchRunner
}

func NewDispatcher(search SearchRunner, providers []provider.Provider) (*Dispatcher, error) {
	d := &Dispatcher{
		providers: make(map[string]provider.Provider),
		search:    search,
	}

	for _, p := range providers {
		for _, name := range p.Names() {
			key := normalize(name)
			if key == "" {
				return nil, fmt.Errorf("provider has empty name")
			}

			if _, exists := d.providers[key]; exists {
				return nil, fmt.Errorf("provider name already registered: %q", key)
			}

			d.providers[key] = p
		}
	}

	return d, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
