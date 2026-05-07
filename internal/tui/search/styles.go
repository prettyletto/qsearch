package search

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/prettyletto/qsearch/internal/domain/provider"
)

type styles struct {
	container  lipgloss.Style
	provider   lipgloss.Style
	input      lipgloss.Style
	selected   lipgloss.Style
	suggestion lipgloss.Style
	hint       lipgloss.Style
}

func newStyles() styles {
	return styles{
		container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6C7086")).
			Padding(1, 2).
			Width(72),
		provider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89B4FA")).
			Bold(true),
		input: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CDD6F4")),
		selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#11111B")).
			Background(lipgloss.Color("#A6E3A1")).
			Bold(true).
			Padding(0, 1).
			Width(66),
		suggestion: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BAC2DE")).
			Padding(0, 1).
			Width(66),
		hint: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086")),
	}
}

func providerLabel(p provider.Provider) string {
	switch p.Names()[0] {
	case "google":
		return "󰊭 google"
	case "youtube":
		return " youtube"
	case "ytmusic":
		return "󰝚 YT Music"
	default:
		return p.Names()[0]
	}
}
