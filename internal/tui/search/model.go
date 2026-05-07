package search

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prettyletto/qsearch/internal/domain/provider"
)

type Result struct {
	Query    string
	Provider provider.Provider
	Canceled bool
}

type model struct {
	styles      styles
	input       textinput.Model
	providers   []provider.Provider
	provider    provider.Provider
	suggestions []string
	selected    int
	result      Result
	width       int
	height      int
}

type suggestionsMsg struct {
	query       string
	suggestions []string
	err         error
}

func New(providers []provider.Provider, provider provider.Provider) model {
	input := textinput.New()
	input.Prompt = " "
	input.Placeholder = "Search " + provider.Names()[0]
	input.Focus()
	input.CharLimit = 200
	input.Width = 60

	return model{
		input:     input,
		providers: providers,
		provider:  provider,
		styles:    newStyles(),
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = m.contentWidth() - lipgloss.Width(m.styles.providerTag(m.provider))
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

	contentWidth := m.contentWidth()
	label := m.styles.providerTag(m.provider)

	input := m.input
	input.Width = max(contentWidth-lipgloss.Width(label)-3, 10)

	inputView := m.styles.input.Render(input.View())
	inputRow := m.styles.inputRow.Render(label + " " + inputView)

	b.WriteString(inputRow)

	maxSuggestions := min(len(m.suggestions), 8)
	var list strings.Builder

	for i := range maxSuggestions {
		suggestion := m.suggestions[i]

		if i == m.selected {
			list.WriteString(m.styles.selectedFor(m.provider, contentWidth).Render(suggestion))
		} else {
			list.WriteString(m.styles.suggestionFor(contentWidth).Render(suggestion))
		}

		list.WriteString("\n")
	}

	b.WriteString(m.styles.list.Render(list.String()))
	b.WriteString(m.styles.footer.Render(m.styles.footerBar()))

	return "\n" + m.styles.containerFor(m.panelWidth()).Render(b.String()) + "\n"
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

func (m model) panelWidth() int {
	if m.width <= 0 {
		return 72
	}

	width := m.width - 4
	if width > 92 {
		return 92
	}
	if width < 56 {
		return 56
	}

	return width
}

func (m model) contentWidth() int {
	return m.panelWidth() - 6
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
