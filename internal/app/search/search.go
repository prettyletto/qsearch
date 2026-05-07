package search

import (
	"strings"

	"github.com/prettyletto/qsearch/internal/domain/provider"

	"github.com/prettyletto/qsearch/internal/infra/browser"
	tuiSearch "github.com/prettyletto/qsearch/internal/tui/search"
)

type Runner struct {
	browser   *browser.Opener
	providers []provider.Provider
}

func NewRunner(browser *browser.Opener, providers []provider.Provider) *Runner {
	return &Runner{
		browser,
		providers,
	}
}

func (r *Runner) Run(p provider.Provider, args []string) error {
	query := strings.TrimSpace(strings.Join(args, " "))

	if query == "" {
		result, err := tuiSearch.Run(r.providers, p)
		if err != nil {
			return err
		}

		if result.Canceled {
			return nil
		}

		query = result.Query
		p = result.Provider
	}

	finalURL := p.SearchURL(query)

	return r.browser.Open(finalURL)
}
