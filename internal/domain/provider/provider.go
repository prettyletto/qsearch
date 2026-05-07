package provider

import "context"

type Provider interface {
	Names() []string
	SearchURL(query string) string
	Suggestions(ctx context.Context, query string) ([]string, error)
}

type Meta struct {
	Icon      string
	Name      string
	TagBG     string
	IconColor string
	TextColor string
}

type MetadataProvider interface {
	Provider
	Meta() Meta
}
