package modes

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jkeresman01/tsm/tmux"
	"github.com/jkeresman01/tsm/utils"
)

type SwitchMode struct {
	sessions []string
	filtered []string
	cursor   int
	input    textinput.Model
}

func NewSwitchMode(sessions []string) *SwitchMode {
	return &SwitchMode{
		sessions: sessions,
		filtered: sessions,
		input:    newSwitchInput(),
	}
}

func (m *SwitchMode) Update(msg tea.Msg) (ModeStrategy, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if next, cmd, done := m.handleKey(key); done {
			return next, cmd
		}
	}
	cmd := m.updateInput(msg)
	m.applyFilter()
	m.clampCursor()
	return m, cmd
}

func (m *SwitchMode) View() string {
	q := m.query()
	var b strings.Builder
	for i, s := range m.filtered {
		b.WriteString(m.rowPrefix(i))
		b.WriteString(utils.HighlightMatches(s, q))
		b.WriteByte('\n')
	}
	return b.String()
}

func (m *SwitchMode) ModeName() string { return "SWITCH MODE" }

func (m *SwitchMode) Reset() { m.input.Reset() }

func (m *SwitchMode) GetCurrentSession() string {
	if m.hasSelection() {
		return m.filtered[m.cursor]
	}
	return ""
}

func newSwitchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search sessions..."
	ti.Focus()
	ti.Prompt = ""
	ti.CharLimit = 64
	ti.Width = 20
	return ti
}

func (m *SwitchMode) handleKey(k tea.KeyMsg) (ModeStrategy, tea.Cmd, bool) {
	switch k.String() {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		if m.hasSelection() {
			tmux.AttachSession(m.filtered[m.cursor])
			return m, tea.Quit, true
		}
	case "d", "delete":
		if m.hasSelection() {
			// Does nothing in this exact moment in the universe
		}
	}
	return nil, nil, false
}

func (m *SwitchMode) updateInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *SwitchMode) applyFilter() {
	m.filtered = utils.FuzzyFilter(m.sessions, m.query())
}

func (m *SwitchMode) clampCursor() {
	n := len(m.filtered)
	if n == 0 {
		m.cursor = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	} else if m.cursor >= n {
		m.cursor = n - 1
	}
}

func (m *SwitchMode) moveCursor(delta int) {
	if len(m.filtered) == 0 {
		m.cursor = 0
		return
	}
	m.cursor += delta
	m.clampCursor()
}

func (m *SwitchMode) hasSelection() bool {
	return len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered)
}

func (m *SwitchMode) query() string {
	return m.input.Value()
}

func (m *SwitchMode) rowPrefix(i int) string {
	if i == m.cursor {
		return "> "
	}
	return "  "
}
