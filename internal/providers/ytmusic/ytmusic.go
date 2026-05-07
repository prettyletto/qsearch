package ytmusic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	searchURL  = "https://music.youtube.com/search"
	suggestURL = "https://music.youtube.com/youtubei/v1/music/get_search_suggestions?alt=json"
)

type suggestionsResponse struct {
	Contents []struct {
		SearchSuggestionsSectionRenderer struct {
			Contents []struct {
				SearchSuggestionRenderer struct {
					NavigationEndpoint struct {
						SearchEndpoint struct {
							Query string `json:"query"`
						} `json:"searchEndpoint"`
					} `json:"navigationEndpoint"`
				} `json:"searchSuggestionRenderer"`
			} `json:"contents"`
		} `json:"searchSuggestionsSectionRenderer"`
	} `json:"contents"`
}

type Provider struct {
	client *http.Client
}

func New() *Provider {
	return &Provider{client: http.DefaultClient}
}

func (p *Provider) Names() []string {
	return []string{"ytmusic", "ym", "music"}
}

func (p *Provider) SearchURL(query string) string {
	u, _ := url.Parse(searchURL)

	values := u.Query()
	values.Set("q", query)

	u.RawQuery = values.Encode()
	return u.String()
}

func (p *Provider) Suggestions(ctx context.Context, query string) ([]string, error) {
	body := map[string]any{
		"input": query,
		"context": map[string]any{
			"client": map[string]any{
				"clientName":    "WEB_REMIX",
				"clientVersion": clientVersion(),
			},
			"user": map[string]any{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, suggestURL, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://music.youtube.com")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("youtube music suggestions returned %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseSuggestions(respBody)
}

func clientVersion() string {
	return "1." + time.Now().UTC().Format("20060102") + ".01.00"
}

func parseSuggestions(body []byte) ([]string, error) {
	var payload suggestionsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	if len(payload.Contents) == 0 {
		return nil, nil
	}

	rawSuggestions := payload.Contents[0].SearchSuggestionsSectionRenderer.Contents
	suggestions := make([]string, 0, len(rawSuggestions))

	for _, raw := range rawSuggestions {
		query := raw.SearchSuggestionRenderer.NavigationEndpoint.SearchEndpoint.Query
		if query != "" {
			suggestions = append(suggestions, query)
		}
	}

	return suggestions, nil
}
