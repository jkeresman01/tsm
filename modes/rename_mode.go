package modes

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jkeresman01/tsm/tmux"
)

type RenameMode struct {
	session string
	input   textinput.Model
}

func NewRenameMode(session string) *RenameMode {
	ti := textinput.New()
	ti.Placeholder = "New name..."
	ti.Prompt = ""
	ti.SetValue(session)
	ti.Focus()
	ti.CharLimit = 64
	return &RenameMode{session, ti}
}

func (m *RenameMode) Update(msg tea.Msg) (ModeStrategy, tea.Cmd) {
	var cmd tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		newName := m.input.Value()
		if newName != "" && newName != m.session {
			tmux.RenameSession(m.session, newName)
		}
		return NewSwitchMode([]string{}), nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *RenameMode) View() string     { return "Rename: " + m.input.View() }
func (m *RenameMode) ModeName() string { return "RENAME MODE" }
func (m *RenameMode) Reset()           { m.input.Reset() }
