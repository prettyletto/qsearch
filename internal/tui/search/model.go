package search

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
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
	keys        keyMap
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

type keyMap struct {
	Provider key.Binding
	Up       key.Binding
	Down     key.Binding
	Open     key.Binding
	Exit     key.Binding

	Google  key.Binding
	YouTube key.Binding
	YTMusic key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Provider: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "provider"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+p"),
			key.WithHelp("^p", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+n"),
			key.WithHelp("^n", "down"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open"),
		),
		Exit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "exit"),
		),
		Google: key.NewBinding(
			key.WithKeys("ctrl+g"),
			key.WithHelp("^g", "google"),
		),
		YouTube: key.NewBinding(
			key.WithKeys("ctrl+y"),
			key.WithHelp("^y", "youtube"),
		),
		YTMusic: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("^u", "ytmusic"),
		),
	}
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
		keys:      newKeyMap(),
	}
}

func Run(providers []provider.Provider, active provider.Provider) (Result, error) {
	m := New(providers, active)

	finalModel, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
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
	return tea.Batch(tea.ClearScreen, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			m.result.Canceled = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Down):
			if len(m.suggestions) > 0 && m.selected < len(m.suggestions)-1 {
				m.selected++
			}
			return m, nil
		case key.Matches(msg, m.keys.Up):
			if len(m.suggestions) > 0 && m.selected > 0 {
				m.selected--
			}
			return m, nil
		case key.Matches(msg, m.keys.Provider):
			return m.cycleProvider()
		case key.Matches(msg, m.keys.Google):
			return m.switchProviderByName("google")
		case key.Matches(msg, m.keys.YouTube):
			return m.switchProviderByName("youtube")
		case key.Matches(msg, m.keys.YTMusic):
			return m.switchProviderByName("ytmusic")
		case key.Matches(msg, m.keys.Open):
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
		m.input.Width = max(m.contentWidth()-lipgloss.Width(m.styles.providerTag(m.provider))-3,
			10)
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
	contentHeight := m.contentHeight()

	label := m.styles.providerTag(m.provider)

	input := m.input
	input.Width = max(contentWidth-lipgloss.Width(label)-3, 10)

	inputView := m.styles.input.Render(input.View())
	inputRow := m.styles.inputRow.Render(label + " " + inputView)

	footerView := m.styles.footer.
		Width(contentWidth).
		Render(m.styles.footerBar(contentWidth, m.footerHints(contentWidth)))

	usedHeight := lipgloss.Height(inputRow) + lipgloss.Height(footerView)
	listHeight := contentHeight - usedHeight - 2
	if listHeight < 0 {
		listHeight = 0
	}

	visibleSuggestions := suggestionLimitForHeight(contentHeight, listHeight)

	visibleStart := 0
	visibleEnd := min(len(m.suggestions), visibleSuggestions)

	if len(m.suggestions) > visibleSuggestions && visibleSuggestions > 0 {
		visibleStart = m.selected - visibleSuggestions/2

		if visibleStart < 0 {
			visibleStart = 0
		}

		maxStart := len(m.suggestions) - visibleSuggestions
		if visibleStart > maxStart {
			visibleStart = maxStart
		}

		visibleEnd = visibleStart + visibleSuggestions
	}

	var list strings.Builder
	for i := visibleStart; i < visibleEnd; i++ {
		suggestion := m.suggestions[i]

		if i == m.selected {
			list.WriteString(m.styles.selectedFor(m.provider, contentWidth).Render(suggestion))
		} else {
			list.WriteString(m.styles.suggestionFor(contentWidth).Render(suggestion))
		}

		list.WriteString("\n")
	}

	listView := m.styles.list.Render(list.String())

	emptyLines := contentHeight -
		lipgloss.Height(inputRow) -
		lipgloss.Height(listView) -
		lipgloss.Height(footerView)
	if emptyLines > 0 {
		listView += strings.Repeat("\n", emptyLines)
	}

	b.WriteString(inputRow)
	b.WriteString(listView)
	b.WriteString(footerView)
	b.WriteString("\n")

	return m.styles.containerFor(m.appWidth()).Render(b.String())
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

func (m model) appWidth() int {
	if m.width <= 0 {
		return 72
	}

	width := m.width
	if width > 120 {
		width = 120
	}
	if width < 56 {
		width = 56
	}

	return width
}

func suggestionLimitForHeight(contentHeight, availableListHeight int) int {
	if availableListHeight <= 0 {
		return 0
	}

	limit := availableListHeight

	if contentHeight <= 8 {
		limit = min(limit, 3)
	} else if contentHeight <= 12 {
		limit = min(limit, 5)
	}

	return limit
}

func (m model) appHeight() int {
	if m.height <= 0 {
		return 18
	}

	if m.height <= 10 {
		return 10
	}

	return m.height
}

func (m model) contentWidth() int {
	return m.appWidth() - 4
}

func (m model) contentHeight() int {
	return m.appHeight() - 1
}

func (m model) footerHints(width int) []footerHint {
	if width < 64 {
		return []footerHint{
			{binding: m.keys.Provider},
			{binding: m.keys.Exit},
		}
	}

	return []footerHint{
		{binding: m.keys.Provider},
		{key: "↑/↓", label: "select"},
		{binding: m.keys.Open},
		{binding: m.keys.Exit},
	}
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
