package modes

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jkeresman01/tsm/tmux"
	"github.com/jkeresman01/tsm/utils"
)

type CreateMode struct {
	dirs     []string
	filtered []string
	cursor   int
	input    textinput.Model
}

func NewCreateMode(dirs []string) *CreateMode {
	return &CreateMode{
		dirs:     dirs,
		filtered: dirs,
		input:    newSearchInput(),
	}
}

func (m *CreateMode) Update(msg tea.Msg) (ModeStrategy, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if next, done := m.handleKey(key); done {
			return next, nil
		}
	}
	cmd := m.updateQuery(msg)
	m.applyFilter()
	m.clampCursor()
	return m, cmd
}

func (m *CreateMode) View() string {
	q := m.query()
	var b strings.Builder

	// Add search bar at top
	b.WriteString("🔍 " + m.input.View() + "\n\n")

	if len(m.filtered) == 0 {
		b.WriteString("  No directories found\n")
		b.WriteString("  Tip: Add search paths in ~/.config/tsm/config.json\n")
		return b.String()
	}

	for i, d := range m.filtered {
		icon := "󰉋 " // Folder icon
		prefix := m.rowPrefix(i)

		b.WriteString(prefix)
		b.WriteString(icon)
		b.WriteString(utils.HighlightMatches(filepath.Base(d), q))

		// Show path hint for selected item
		if i == m.cursor {
			b.WriteString("  󰄾") // Arrow indicator
		}

		b.WriteByte('\n')

		// Show full path for current selection
		if i == m.cursor {
			pathStyle := "    󰉖 " + d
			b.WriteString(pathStyle)
			b.WriteByte('\n')
		}
	}

	// Add count at bottom
	b.WriteString(fmt.Sprintf("\n  %d director(ies)", len(m.filtered)))

	return b.String()
}

func (m *CreateMode) ModeName() string          { return "CREATE" }
func (m *CreateMode) Reset()                    {}
func (m *CreateMode) GetCurrentSession() string { return "" }

func newSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search directories..."
	ti.Focus()
	ti.Prompt = ""
	ti.Width = 30
	return ti
}

func (m *CreateMode) handleKey(k tea.KeyMsg) (ModeStrategy, bool) {
	switch k.String() {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		return m.onConfirm(), true
	}
	return nil, false
}

func (m *CreateMode) updateQuery(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *CreateMode) applyFilter() {
	m.filtered = utils.FuzzyFilter(m.dirs, m.query())
}

func (m *CreateMode) moveCursor(delta int) {
	if len(m.filtered) == 0 {
		m.cursor = 0
		return
	}
	m.cursor += delta
	m.clampCursor()
}

func (m *CreateMode) clampCursor() {
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

func (m *CreateMode) onConfirm() ModeStrategy {
	if !m.hasSelection() {
		return m
	}
	dir := m.selectedDir()
	name := filepath.Base(dir)
	tmux.CreateSession(name, dir)

	sessions, _ := tmux.ListSessions()
	return NewSwitchMode(sessions)
}

func (m *CreateMode) hasSelection() bool {
	return len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered)
}

func (m *CreateMode) selectedDir() string {
	return m.filtered[m.cursor]
}

func (m *CreateMode) rowPrefix(i int) string {
	if i == m.cursor {
		return "▶ "
	}
	return "  "
}

func (m *CreateMode) query() string {
	return m.input.Value()
}
