package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	modes "github.com/jkeresman01/tsm/modes"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/jkeresman01/tsm/tmux"
)

type SessionModel struct {
	sessions     []string
	filtered     []string
	cursor       int
	showInput    bool
	inputFocused bool
	Mode         modes.ModeStrategy
	input        textinput.Model
}

func (m *SessionModel) SetMode(mode modes.ModeStrategy) {
	m.Mode = mode
	m.Mode.Reset()
}

func (m *SessionModel) NextMode() {
	switch m.Mode.(type) {
	case *modes.SwitchMode:
		m.Mode = modes.NewRenameMode("default-session")
	case *modes.RenameMode:
		m.Mode = modes.NewCreateMode([]string{})
	case *modes.CreateMode:
		m.Mode = modes.NewSwitchMode([]string{})
	}
}

func (m SessionModel) Input() textinput.Model {
	return m.input
}

func (m SessionModel) Sessions() []string {
	return m.sessions
}

func (m SessionModel) CurrentSession() string {
	if len(m.filtered) == 0 {
		return ""
	}
	return m.filtered[m.cursor]
}

func NewSessionModel() SessionModel {
	sessions, err := tmux.ListSessions()
	if err != nil || len(sessions) == 0 {
		sessions = []string{"No sessions"}
	}

	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	ti.Prompt = ""
	ti.CharLimit = 64
	ti.Width = 20

	return SessionModel{
		sessions:  sessions,
		filtered:  sessions,
		showInput: true,
		input:     ti,
	}
}

func (m SessionModel) Update(msg tea.Msg) (SessionModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter":
			// TODO: attach to session
		}
	}

	m.input, cmd = m.input.Update(msg)
	query := m.input.Value()
	m.filtered = fuzzyFilter(m.sessions, query)
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}

	return m, cmd
}

func (m SessionModel) View() string {
	query := m.input.Value()
	var s string

	for i, sess := range m.filtered {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		display := highlightMatch(sess, query)
		s += cursor + " " + display + "\n"
	}
	return s
}

func fuzzyFilter(items []string, query string) []string {
	if query == "" {
		return items
	}
	var out []string
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), query) {
			out = append(out, item)
		}
	}
	return out
}

var matchStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("212"))

func highlightMatch(session, query string) string {
	if query == "" {
		return session
	}
	q := strings.ToLower(query)
	s := strings.ToLower(session)

	var result strings.Builder

	for i := 0; i < len(session); i++ {
		c := session[i]
		if strings.ContainsRune(q, rune(s[i])) {
			result.WriteString(matchStyle.Render(string(c)))
		} else {
			result.WriteByte(c)
		}
	}
	return result.String()
}
