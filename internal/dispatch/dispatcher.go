package dispatch

import (
	"fmt"
	"sort"
	"strings"

	"github.com/prettyletto/qseach/internal/domain/provider"
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

func (d *Dispatcher) Dispatch(args []string) error {
	if len(args) == 0 {

		return d.printHelp()
	}

	first := normalize(args[0])

	switch first {
	case "help", "-h", "--help":
		return d.printHelp()

	}

	p, ok := d.providers[first]
	if !ok {
		return fmt.Errorf("unkown provider %q", args[0])
	}

	return d.search.Run(p, args[1:])
}

func (d *Dispatcher) printHelp() error {
	fmt.Println("usage: qs <provider> [query]")
	fmt.Println()
	fmt.Println("providers:")

	names := d.providersNames()
	for _, name := range names {
		fmt.Printf("  %s\n", name)
	}

	return nil
}

func (d *Dispatcher) providersNames() []string {
	seen := make(map[provider.Provider]bool)
	names := make([]string, 0)

	for _, p := range d.providers {
		if seen[p] {
			continue
		}

		names = append(names, p.Names()[0])

		seen[p] = true
	}

	sort.Strings(names)
	return names
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
