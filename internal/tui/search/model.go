package search

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prettyletto/qsearch/internal/domain/provider"
)

type Result struct {
	Query    string
	Provider provider.Provider
	Canceled bool
}

type model struct {
	input       textinput.Model
	providers   []provider.Provider
	provider    provider.Provider
	suggestions []string
	selected    int
	result      Result
}

type suggestionsMsg struct {
	query       string
	suggestions []string
	err         error
}

func New(providers []provider.Provider, provider provider.Provider) model {
	input := textinput.New()
	input.Placeholder = "Search " + provider.Names()[0]
	input.Focus()
	input.CharLimit = 200
	input.Width = 60

	return model{
		input:     input,
		providers: providers,
		provider:  provider,
	}
}

func Run(providers []provider.Provider, active provider.Provider) (Result, error) {
	m := New(providers, active)

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return Result{}, err
	}

	m, ok := finalModel.(model)
	if !ok {
		return Result{}, nil
	}

	return m.result, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.result.Canceled = true
			return m, tea.Quit
		case "down", "ctrl+n":
			if len(m.suggestions) > 0 && m.selected < len(m.suggestions)-1 {
				m.selected++
			}
			return m, nil
		case "up", "ctrl+p":
			if len(m.suggestions) > 0 && m.selected > 0 {
				m.selected--
			}
			return m, nil
		case "tab":
			return m.cycleProvider()
		case "ctrl+g":
			return m.switchProviderByName("google")
		case "ctrl+y":
			return m.switchProviderByName("youtube")
		case "ctrl+u":
			return m.switchProviderByName("ytmusic")
		case "enter":
			query := strings.TrimSpace(m.input.Value())

			if len(m.suggestions) > 0 && m.selected >= 0 {
				query = m.suggestions[m.selected]
			}

			if query == "" {
				return m, nil
			}

			m.result.Query = query
			m.result.Provider = m.provider
			return m, tea.Quit
		}
	case suggestionsMsg:
		if msg.err != nil {
			return m, nil
		}

		if strings.TrimSpace(m.input.Value()) != msg.query {
			return m, nil
		}

		m.suggestions = msg.suggestions
		m.selected = 0
		return m, nil
	}

	oldValue := m.input.Value()
	m.input, cmd = m.input.Update(msg)

	newValue := strings.TrimSpace(m.input.Value())
	if newValue != strings.TrimSpace(oldValue) {
		if newValue == "" {
			m.suggestions = nil
			m.selected = 0
			return m, cmd
		}

		return m, tea.Batch(cmd, fetchSuggestions(m.provider, newValue))
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n ")
	b.WriteString("[")
	b.WriteString(m.provider.Names()[0])
	b.WriteString("]: ")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	for i, suggestion := range m.suggestions {
		cursor := " "
		if i == m.selected {
			cursor = "> "
		}
		b.WriteString("  ")
		b.WriteString(cursor)
		b.WriteString(suggestion)
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) switchProvider(p provider.Provider) (model, tea.Cmd) {
	m.provider = p
	m.suggestions = nil
	m.selected = 0

	query := strings.TrimSpace(m.input.Value())
	if query == "" {
		return m, nil
	}

	return m, fetchSuggestions(m.provider, query)
}

func (m model) switchProviderByName(name string) (model, tea.Cmd) {
	idx := slices.IndexFunc(m.providers, func(p provider.Provider) bool {
		return p.Names()[0] == name
	})
	if idx == -1 {
		return m, nil
	}

	return m.switchProvider(m.providers[idx])
}

func (m model) cycleProvider() (model, tea.Cmd) {
	idx := slices.IndexFunc(m.providers, func(p provider.Provider) bool {
		return p.Names()[0] == m.provider.Names()[0]
	})
	if idx == -1 {
		return m, nil
	}

	idx = (idx + 1) % len(m.providers)

	return m.switchProvider(m.providers[idx])
}

func fetchSuggestions(p provider.Provider, query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		suggestions, err := p.Suggestions(ctx, query)

		return suggestionsMsg{
			query,
			suggestions,
			err,
		}
	}
}
