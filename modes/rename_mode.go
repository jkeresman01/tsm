package modes

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jkeresman01/tsm/tmux"
	"github.com/jkeresman01/tsm/utils"
)

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			RenameMode handles renaming existing tmux sessions.
//
//		@Description	Provides session selection and renaming interface
//
// ///////////////////////////////////////////////////////////////////////////////////////////
type RenameMode struct {
	sessions        []string
	filtered        []string
	cursor          int
	searchInput     textinput.Model
	renameInput     textinput.Model
	renaming        bool
	selectedSession string
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			NewRenameMode creates a new RenameMode instance.
//
//		@Param			session	string		Session to rename (empty for selection mode)
//
//		@Return			*RenameMode	Initialized RenameMode
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			Update processes input messages and updates the mode state.
//
//		@Param			msg		tea.Msg			Input message
//
//		@Return			ModeStrategy	Updated mode state
//		@Return			tea.Cmd			Optional command
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			View renders the mode's UI.
//
//		@Return			string	Rendered view
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) View() string {
	if m.renaming {
		return m.renderRenamePrompt()
	}
	return m.renderSessionList()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			ModeName returns the display name of this mode.
//
//		@Return			string	Mode name with context
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) ModeName() string {
	if m.renaming {
		return "RENAME: " + m.selectedSession
	}
	return "RENAME MODE - SELECT SESSION"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	@Brief			Reset clears the mode state.
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) Reset() {
	m.searchInput.Reset()
	m.renameInput.Reset()
	m.renaming = false
	m.selectedSession = ""
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			GetCurrentSession returns the currently selected or renaming session.
//
//		@Return			string	Session name or empty string
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) GetCurrentSession() string {
	if m.renaming {
		return m.selectedSession
	}
	if m.hasSelection() {
		return m.filtered[m.cursor]
	}
	return ""
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			handleKey processes keyboard input.
//
//		@Param			k		tea.KeyMsg		Keyboard message
//
//		@Return			ModeStrategy	Next mode (if changed)
//		@Return			tea.Cmd			Command to execute
//		@Return			bool			Whether key was handled
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) handleKey(k tea.KeyMsg) (ModeStrategy, tea.Cmd, bool) {
	if m.renaming {
		return m.handleRenameKeys(k)
	}
	return m.handleSelectionKeys(k)
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			handleRenameKeys processes keys during rename input.
//
//		@Param			k		tea.KeyMsg		Keyboard message
//
//		@Return			ModeStrategy	Next mode (if changed)
//		@Return			tea.Cmd			Command to execute
//		@Return			bool			Whether key was handled
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			handleSelectionKeys processes keys during session selection.
//
//		@Param			k		tea.KeyMsg		Keyboard message
//
//		@Return			ModeStrategy	Next mode (if changed)
//		@Return			tea.Cmd			Command to execute
//		@Return			bool			Whether key was handled
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			GetIcon returns the mode's icon.
//
//		@Return			string	Nerd font icon
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) GetIcon() string {
	return "󰑕"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			GetFooterText returns the help text for the footer.
//
//		@Return			string	Keybinding help text
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) GetFooterText() string {
	if m.renaming {
		return "type new name • ↵ confirm • ⎋ cancel • ? help • q quit"
	}
	return "↑↓ navigate • ↵ select/confirm • ⎋ cancel • ⇥ cycle • ? help • q quit"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			updateRenameInput updates the rename input field.
//
//		@Param			msg		tea.Msg		Input message
//
//		@Return			ModeStrategy	This mode
//		@Return			tea.Cmd			Command from input update
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) updateRenameInput(msg tea.Msg) (ModeStrategy, tea.Cmd) {
	var cmd tea.Cmd
	m.renameInput, cmd = m.renameInput.Update(msg)
	return m, cmd
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			updateSearch updates the search input field.
//
//		@Param			msg		tea.Msg		Input message
//
//		@Return			tea.Cmd	Command from input update
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) updateSearch(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return cmd
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	@Brief			applyFilter filters sessions based on search query.
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) applyFilter() {
	m.filtered = utils.FuzzyFilter(m.sessions, m.query())
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			moveCursor moves the cursor by delta positions.
//
//		@Param			delta	int	Number of positions to move (negative for up)
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) moveCursor(delta int) {
	if len(m.filtered) == 0 {
		m.cursor = 0
		return
	}
	m.cursor += delta
	m.clampCursor()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	@Brief			clampCursor ensures cursor stays within valid range.
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			confirmRename performs the rename operation.
//
//		@Return			ModeStrategy	SwitchMode with updated session list
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) confirmRename() ModeStrategy {
	newName := m.renameInput.Value()
	if newName != "" && newName != m.selectedSession {
		tmux.RenameSession(m.selectedSession, newName)
	}
	sessions, _ := tmux.ListSessions()
	return NewSwitchMode(sessions)
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	@Brief			cancelRename cancels the rename operation.
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) cancelRename() {
	m.renaming = false
	m.renameInput.Reset()
	m.searchInput.Focus()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	@Brief			startRename begins the rename process for selected session.
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			hasSelection returns whether a valid session is selected.
//
//		@Return			bool	True if a session is selected
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) hasSelection() bool {
	return len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered)
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			query returns the current search query.
//
//		@Return			string	Search query text
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) query() string {
	return m.searchInput.Value()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			rowPrefix returns the prefix for a list row.
//
//		@Param			i		int		Row index
//
//		@Return			string	Prefix ("> " for selected, "  " for others)
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) rowPrefix(i int) string {
	if i == m.cursor {
		return "> "
	}
	return "  "
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			renderRenamePrompt renders the rename input prompt.
//
//		@Return			string	Rendered prompt
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func (m *RenameMode) renderRenamePrompt() string {
	return "Rename '" + m.selectedSession + "' to: " + m.renameInput.View()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			renderSessionList renders the session selection list.
//
//		@Return			string	Rendered session list
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			newRenameSearchInput creates a configured search input.
//
//		@Return			textinput.Model	Configured input field
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func newRenameSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search sessions..."
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 20
	return ti
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			newRenameInput creates a configured rename input.
//
//		@Return			textinput.Model	Configured input field
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func newRenameInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "New name..."
	ti.Prompt = ""
	ti.CharLimit = 64
	return ti
}
