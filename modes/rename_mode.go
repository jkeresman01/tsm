package modes

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jkeresman01/tsm/tmux"
	"github.com/jkeresman01/tsm/utils"
)

type RenameMode struct {
	sessions        []string
	filtered        []string
	cursor          int
	searchInput     textinput.Model
	renameInput     textinput.Model
	renaming        bool
	selectedSession string
}

func NewRenameMode(session string) *RenameMode {
	sessions, _ := tmux.ListSessions()
	searchInput := newRenameSearchInput()
	renameInput := newRenameInput()
	renaming := session != ""

	if renaming {
		renameInput.SetValue(session)
		renameInput.Focus()
		searchInput.Blur()
	}

	return &RenameMode{
		sessions:        sessions,
		filtered:        sessions,
		searchInput:     searchInput,
		renameInput:     renameInput,
		renaming:        renaming,
		selectedSession: session,
	}
}

func (m *RenameMode) Update(msg tea.Msg) (ModeStrategy, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if next, cmd, done := m.handleKey(key); done {
			return next, cmd
		}
	}

	if m.renaming {
		return m.updateRenameInput(msg)
	}

	cmd := m.updateSearch(msg)
	m.applyFilter()
	m.clampCursor()
	return m, cmd
}

func (m *RenameMode) View() string {
	if m.renaming {
		return m.renderRenamePrompt()
	}
	return m.renderSessionList()
}

func (m *RenameMode) ModeName() string {
	if m.renaming {
		return "RENAME: " + m.selectedSession
	}
	return "RENAME MODE - SELECT SESSION"
}

func (m *RenameMode) Reset() {
	m.searchInput.Reset()
	m.renameInput.Reset()
	m.renaming = false
	m.selectedSession = ""
}

func (m *RenameMode) GetCurrentSession() string {
	if m.renaming {
		return m.selectedSession
	}
	if m.hasSelection() {
		return m.filtered[m.cursor]
	}
	return ""
}

func (m *RenameMode) handleKey(k tea.KeyMsg) (ModeStrategy, tea.Cmd, bool) {
	if m.renaming {
		return m.handleRenameKeys(k)
	}
	return m.handleSelectionKeys(k)
}

func (m *RenameMode) handleRenameKeys(k tea.KeyMsg) (ModeStrategy, tea.Cmd, bool) {
	switch k.String() {
	case "enter":
		return m.confirmRename(), nil, true
	case "esc":
		m.cancelRename()
		return m, nil, true
	}
	return nil, nil, false
}

func (m *RenameMode) handleSelectionKeys(k tea.KeyMsg) (ModeStrategy, tea.Cmd, bool) {
	switch k.String() {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		m.startRename()
		return m, nil, true
	case "esc":
		sessions, _ := tmux.ListSessions()
		return NewSwitchMode(sessions), nil, true
	}
	return nil, nil, false
}

func (m *RenameMode) GetIcon() string {
	return "󰑕"
}

func (m *RenameMode) GetFooterText() string {
	if m.renaming {
		return "type new name • ↵ confirm • ⎋ cancel • ? help • q quit"
	}
	return "↑↓ navigate • ↵ select/confirm • ⎋ cancel • ⇥ cycle • ? help • q quit"
}

func (m *RenameMode) updateRenameInput(msg tea.Msg) (ModeStrategy, tea.Cmd) {
	var cmd tea.Cmd
	m.renameInput, cmd = m.renameInput.Update(msg)
	return m, cmd
}

func (m *RenameMode) updateSearch(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return cmd
}

func (m *RenameMode) applyFilter() {
	m.filtered = utils.FuzzyFilter(m.sessions, m.query())
}

func (m *RenameMode) moveCursor(delta int) {
	if len(m.filtered) == 0 {
		m.cursor = 0
		return
	}
	m.cursor += delta
	m.clampCursor()
}

func (m *RenameMode) clampCursor() {
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

func (m *RenameMode) confirmRename() ModeStrategy {
	newName := m.renameInput.Value()
	if newName != "" && newName != m.selectedSession {
		tmux.RenameSession(m.selectedSession, newName)
	}
	sessions, _ := tmux.ListSessions()
	return NewSwitchMode(sessions)
}

func (m *RenameMode) cancelRename() {
	m.renaming = false
	m.renameInput.Reset()
	m.searchInput.Focus()
}

func (m *RenameMode) startRename() {
	if !m.hasSelection() {
		return
	}
	m.selectedSession = m.filtered[m.cursor]
	m.renameInput.SetValue(m.selectedSession)
	m.renameInput.Focus()
	m.searchInput.Blur()
	m.renaming = true
}

func (m *RenameMode) hasSelection() bool {
	return len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered)
}

func (m *RenameMode) query() string {
	return m.searchInput.Value()
}

func (m *RenameMode) rowPrefix(i int) string {
	if i == m.cursor {
		return "> "
	}
	return "  "
}

func (m *RenameMode) renderRenamePrompt() string {
	return "Rename '" + m.selectedSession + "' to: " + m.renameInput.View()
}

func (m *RenameMode) renderSessionList() string {
	var b strings.Builder
	b.WriteString("Select session to rename:\n\n")
	q := m.query()
	for i, s := range m.filtered {
		b.WriteString(m.rowPrefix(i))
		b.WriteString(utils.HighlightMatches(s, q))
		b.WriteByte('\n')
	}
	return b.String()
}

func newRenameSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search sessions..."
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 20
	return ti
}

func newRenameInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "New name..."
	ti.Prompt = ""
	ti.CharLimit = 64
	return ti
}
