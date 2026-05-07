package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prettyletto/qsearch/internal/domain/provider"

	"github.com/prettyletto/qsearch/internal/infra/browser"
	tuiSearch "github.com/prettyletto/qsearch/internal/tui/search"
)

type Runner struct {
	browser *browser.Opener
}

func NewRunner(browser *browser.Opener) *Runner {
	return &Runner{
		browser,
	}
}

func (r *Runner) Run(p provider.Provider, args []string) error {

	query := strings.TrimSpace(strings.Join(args, " "))

	if query == "" {
		result, err := tuiSearch.Run(p)
		if err != nil {
			return err
		}

		if result.Canceled {
			return nil
		}

		query = result.Query
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	suggestions, err := p.Suggestions(ctx, query)

	if err != nil {
		return err
	}

	if len(suggestions) > 0 {
		for _, suggestion := range suggestions {
			fmt.Println(suggestion)
		}
	}

	finalURL := p.SearchURL(query)

	return r.browser.Open(finalURL)
}
