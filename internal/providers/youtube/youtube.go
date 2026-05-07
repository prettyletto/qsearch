package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	searchURL  = "https://youtube.com/results"
	suggestURL = "https://suggestqueries.google.com/complete/search"
)

type Provider struct {
	client *http.Client
}

func New() *Provider {
	return &Provider{
		client: http.DefaultClient,
	}
}

func (p *Provider) Names() []string {
	return []string{"youtube", "y", "yt"}
}

func (p *Provider) SearchURL(query string) string {
	u, _ := url.Parse(searchURL)

	values := u.Query()
	values.Set("search_query", query)

	u.RawQuery = values.Encode()
	return u.String()
}

func (p *Provider) Suggestions(ctx context.Context, query string) ([]string, error) {
	u, _ := url.Parse(suggestURL)

	values := u.Query()
	values.Set("client", "firefox")
	values.Set("ds", "yt")
	values.Set("q", query)
	values.Set("ie", "utf-8")
	values.Set("oe", "utf-8")

	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("youtube susggestions returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseSuggestions(body)
}

func parseSuggestions(body []byte) ([]string, error) {
	var payload []json.RawMessage

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	if len(payload) < 2 {
		return nil, fmt.Errorf("unexpected youtube suggestions response")
	}

	var suggestions []string
	if err := json.Unmarshal(payload[1], &suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}
