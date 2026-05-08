package search

import (
	"github.com/charmbracelet/bubbles/key"
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

type theme struct {
	text        lipgloss.TerminalColor
	muted       lipgloss.TerminalColor
	subtle      lipgloss.TerminalColor
	panelBorder lipgloss.TerminalColor
	rowBG       lipgloss.TerminalColor
	keyBG       lipgloss.TerminalColor
	tagFallback lipgloss.TerminalColor
}

type providerMeta struct {
	icon string
	name string

	iconColor lipgloss.Color
	textColor lipgloss.Color
	tagBG     lipgloss.Color
}

type footerHint struct {
	binding key.Binding
	key     string
	label   string
}

func defaultTheme() theme {
	return theme{
		text: lipgloss.AdaptiveColor{
			Light: "#242424",
			Dark:  "15",
		},
		muted: lipgloss.Color("8"),
		subtle: lipgloss.AdaptiveColor{
			Light: "#5A5A5A",
			Dark:  "7",
		},
		panelBorder: lipgloss.Color("8"),
		rowBG: lipgloss.AdaptiveColor{
			Light: "#E8E4D8",
			Dark:  "0",
		},
		keyBG:       lipgloss.Color("8"),
		tagFallback: lipgloss.Color("8"),
	}
}

func newStyles() styles {
	t := defaultTheme()

	return styles{
		container: lipgloss.NewStyle().
			PaddingTop(1).
			PaddingLeft(1).
			PaddingRight(1).
			Width(72),
		input: lipgloss.NewStyle().
			Foreground(t.text),
		inputRow: lipgloss.NewStyle().
			MarginBottom(1),
		list: lipgloss.NewStyle().
			MarginTop(1),
		selected: lipgloss.NewStyle().
			Foreground(t.text).
			Background(t.rowBG).
			Bold(true).
			Padding(0, 1).
			Width(64),
		suggestion: lipgloss.NewStyle().
			Foreground(t.subtle).
			Padding(0, 2).
			Width(66),
		footer: lipgloss.NewStyle().
			Foreground(t.muted).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(t.muted),
		keycap: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("238")).
			Bold(true).
			Padding(0, 1),
		hintText: lipgloss.NewStyle().
			Foreground(t.muted),
		tag: lipgloss.NewStyle().
			Background(t.tagFallback).
			Padding(0, 1),
		tagIcon: lipgloss.NewStyle().
			Bold(true),
		tagText: lipgloss.NewStyle().
			Bold(true),
		tagSuffix: lipgloss.NewStyle().
			Foreground(t.muted),
		tagLeftCap:  lipgloss.NewStyle(),
		tagRightCap: lipgloss.NewStyle(),
	}
}

func (s styles) footerBar(width int, hints []footerHint) string {
	spacer := s.hintText.Render("  ")
	items := make([]string, 0, len(hints))

	for _, hint := range hints {
		keyText, label := hint.key, hint.label
		if keyText == "" || label == "" {
			help := hint.binding.Help()
			if keyText == "" {
				keyText = help.Key
			}
			if label == "" {
				label = help.Desc
			}
		}

		items = append(items, s.keycap.Render(keyText)+s.hintText.Render(" "+label))
	}

	if len(items) == 0 {
		return ""
	}

	line := items[0]
	for _, item := range items[1:] {
		line += spacer + item
	}

	return line
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
		Width(width)
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
