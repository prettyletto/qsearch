package custom

import (
	"context"
	"net/url"
	"strings"

	"github.com/prettyletto/qsearch/internal/domain/provider"
)

type Provider struct {
	names       []string
	urlTemplate string
	meta        provider.Meta
}

func New(names []string, urlTemplate string, meta provider.Meta) *Provider {
	return &Provider{
		names:       names,
		urlTemplate: urlTemplate,
		meta:        meta,
	}
}

func (p *Provider) Names() []string {
	return p.names
}

func (p *Provider) SearchURL(query string) string {
	escapedQuery := url.QueryEscape(query)

	return strings.ReplaceAll(p.urlTemplate, "{{query}}", escapedQuery)
}

// IN LATER VERSIONS THIS SHOULD BE WIRED TOO
func (p *Provider) Suggestions(ctx context.Context, query string) ([]string, error) {
	return nil, nil
}

func (p *Provider) Meta() provider.Meta {
	return p.meta
}
