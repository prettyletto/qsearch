package provider

import "context"

type Provider interface {
	Names() []string
	SearchURL(query string) string
	Suggestions(ctx context.Context, query string) ([]string, error)
}
