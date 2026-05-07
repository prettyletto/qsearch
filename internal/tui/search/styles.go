package search

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/prettyletto/qsearch/internal/domain/provider"
)

type styles struct {
	container  lipgloss.Style
	input      lipgloss.Style
	inputRow   lipgloss.Style
	list       lipgloss.Style
	suggestion lipgloss.Style
	footer     lipgloss.Style
	keycap     lipgloss.Style
	hintText   lipgloss.Style

	tag         lipgloss.Style
	tagIcon     lipgloss.Style
	tagText     lipgloss.Style
	tagSuffix   lipgloss.Style
	tagLeftCap  lipgloss.Style
	tagRightCap lipgloss.Style

	selected lipgloss.Style
}

type providerMeta struct {
	icon string
	name string

	iconColor lipgloss.Color
	textColor lipgloss.Color
	tagBG     lipgloss.Color
}

func newStyles() styles {
	return styles{
		container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6C7086")).
			Padding(1, 2).
			Width(72),
		input: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CDD6F4")),
		inputRow: lipgloss.NewStyle().
			MarginBottom(1),
		list: lipgloss.NewStyle().
			MarginTop(1).
			MarginBottom(1),
		selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#313244")).
			Bold(true).
			Padding(0, 1).
			Width(64),
		suggestion: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BAC2DE")).
			Padding(0, 2).
			Width(66),
		footer: lipgloss.NewStyle().
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#313244")).
			PaddingTop(1),
		keycap: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CDD6F4")).
			Background(lipgloss.Color("#313244")).
			Bold(true).
			Padding(0, 1),
		hintText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086")),
		tag: lipgloss.NewStyle().
			Background(lipgloss.Color("#313244")).
			Padding(0, 1),
		tagIcon: lipgloss.NewStyle().
			Bold(true),
		tagText: lipgloss.NewStyle().
			Bold(true),
		tagSuffix: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086")),
		tagLeftCap:  lipgloss.NewStyle(),
		tagRightCap: lipgloss.NewStyle(),
	}
}

func (s styles) footerBar() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		s.keycap.Render("tab"),
		s.hintText.Render(" switch  "),
		s.keycap.Render("ctrl+g"),
		s.hintText.Render(" google  "),
		s.keycap.Render("ctrl+y"),
		s.hintText.Render(" youtube  "),
		s.keycap.Render("ctrl+u"),
		s.hintText.Render(" music  "),
		s.keycap.Render("esc"),
		s.hintText.Render(" close"),
	)
}

func (s styles) providerTag(p provider.Provider) string {
	meta := providerMetaFor(p)

	leftCap := s.tagLeftCap.
		Foreground(meta.tagBG).
		Render("")

	icon := s.tag.
		Foreground(meta.iconColor).
		Background(meta.tagBG).
		Bold(true).
		Render(meta.icon)

	text := s.tag.
		Foreground(meta.textColor).
		Background(meta.tagBG).
		Bold(true).
		Render(meta.name + " ")

	rightCap := s.tagRightCap.
		Foreground(meta.tagBG).
		Render("")

	return leftCap + icon + text + rightCap + s.tagSuffix.Render(":")
}

func (s styles) containerFor(width int) lipgloss.Style {
	return s.container.
		Width(width).
		BorderForeground(lipgloss.Color("#45475A"))
}

func (s styles) selectedFor(p provider.Provider, width int) lipgloss.Style {
	meta := providerMetaFor(p)

	return s.selected.
		Width(width - 1).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(meta.tagBG)
}

func (s styles) suggestionFor(width int) lipgloss.Style {
	return s.suggestion.Width(width)
}

func providerMetaFor(p provider.Provider) providerMeta {
	if p, ok := p.(provider.MetadataProvider); ok {
		meta := p.Meta()

		return providerMeta{
			icon:      meta.Icon,
			name:      meta.Name,
			iconColor: lipgloss.Color(meta.IconColor),
			textColor: lipgloss.Color(meta.TextColor),
			tagBG:     lipgloss.Color(meta.TagBG),
		}
	}

	switch p.Names()[0] {
	case "google":
		return providerMeta{
			icon:      "󰊭",
			name:      "Google",
			iconColor: lipgloss.Color("#4285F4"),
			textColor: lipgloss.Color("#FFFFFF"),
			tagBG:     lipgloss.Color("#1F1F1F"),
		}
	case "youtube":
		return providerMeta{
			icon:      "",
			name:      "YouTube",
			iconColor: lipgloss.Color("#FFFFFF"),
			textColor: lipgloss.Color("#FFFFFF"),
			tagBG:     lipgloss.Color("#FF0635"),
		}
	case "ytmusic":
		return providerMeta{
			icon:      "",
			name:      "YT Music",
			iconColor: lipgloss.Color("#FF0808"),
			textColor: lipgloss.Color("#2D2D2D"),
			tagBG:     lipgloss.Color("#FFFFFF"),
		}
	default:
		return providerMeta{
			name:      p.Names()[0],
			iconColor: lipgloss.Color("#FFFFFF"),
			textColor: lipgloss.Color("#FFFFFF"),
			tagBG:     lipgloss.Color("#313244"),
		}

	}
}
