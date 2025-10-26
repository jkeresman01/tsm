package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jkeresman01/tsm/config"
	modes "github.com/jkeresman01/tsm/modes"
	styles "github.com/jkeresman01/tsm/styles"
	"github.com/jkeresman01/tsm/tmux"
	"github.com/jkeresman01/tsm/utils"
)

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			manager is the main application model for Bubble Tea.
//
//	@Description	Manages the overall UI state, mode coordination, and layout rendering
//
/////////////////////////////////////////////////////////////////////////////////////////////
type manager struct {
	width    int                // Terminal width
	height   int                // Terminal height
	showHelp bool               // Whether help dialog is visible
	mode     modes.ModeStrategy // Current operational mode
	dirs     []string           // Available project directories
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			NewTsmManager creates a new TSM manager instance.
//
//	@Param			cfg		config.Config	Application configuration
//
//	@Return			tea.Model	Initialized Bubble Tea model
//
/////////////////////////////////////////////////////////////////////////////////////////////
func NewTsmManager(cfg config.Config) tea.Model {
	sessions, _ := tmux.ListSessions()
	if len(sessions) == 0 {
		sessions = []string{}
	}
	dirs := utils.GetProjectDirs(cfg.SearchPaths, cfg.MaxDepth)
	return &manager{
		mode: modes.NewSwitchMode(sessions),
		dirs: dirs,
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			Init initializes the manager (Bubble Tea Init method).
//
//	@Return			tea.Cmd	Initial command (nil in this case)
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) Init() tea.Cmd { return nil }

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			Update processes messages and updates the manager state.
//
//	@Param			msg		tea.Msg		Input message
//
//	@Return			tea.Model	Updated model
//	@Return			tea.Cmd		Command to execute
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.WindowSizeMsg:
		m.applyWindowSize(t)
		return m, nil
	case tea.KeyMsg:
		if cmd := m.handleGlobalKey(t); cmd != nil {
			return m, cmd
		}
	}
	newMode, cmd := m.mode.Update(msg)
	m.mode = newMode
	return m, cmd
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			View renders the complete UI (Bubble Tea View method).
//
//	@Return			string	Rendered view
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.showHelp {
		return m.renderHelpOverlay()
	}
	return m.renderLayout()
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderLayout renders the main application layout.
//
//	@Return			string	Complete rendered layout with header, body, and footer
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderLayout() string {
	header := m.renderHeader()
	body := m.renderBody()
	footer := m.renderFooter()
	padding := strings.Repeat("\n", m.remainingHeight(lipgloss.Height(header)+lipgloss.Height(body)+lipgloss.Height(footer)))
	layout := lipgloss.JoinVertical(lipgloss.Top, header, body, padding, footer)
	withOuter := styles.CurrentTheme.OuterStyle.Render(layout)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, withOuter)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			applyWindowSize updates the manager's width and height.
//
//	@Param			msg		tea.WindowSizeMsg	Window size message
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) applyWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			handleGlobalKey processes global keyboard shortcuts.
//
//	@Param			k		tea.KeyMsg	Keyboard message
//
//	@Return			tea.Cmd	Command to execute (e.g., tea.Quit)
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) handleGlobalKey(k tea.KeyMsg) tea.Cmd {
	switch k.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "?":
		m.showHelp = !m.showHelp
	case "tab":
		m.cycleMode()
	case "ctrl+n":
		m.handleCreateMode()
	case "ctrl+r":
		m.mode = modes.NewRenameMode("")
	case "ctrl+s":
		m.handleSwitchMode()
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			handleCreateMode switches to create mode.
//
//	@Description	Uses existing dirs or loads default directories if none available
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) handleCreateMode() {
	m.mode = modes.NewCreateMode(m.dirs)
	if len(m.dirs) == 0 {
		m.dirs = m.getDefaultDirs()
		m.mode = modes.NewCreateMode(m.dirs)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			handleSwitchMode switches to switch mode with current sessions.
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) handleSwitchMode() {
	sessions, _ := tmux.ListSessions()
	m.mode = modes.NewSwitchMode(sessions)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			cycleMode cycles through the available modes.
//
//	@Description	Order: Switch -> Rename -> Create -> Switch
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) cycleMode() {
	sessions, _ := tmux.ListSessions()
	switch m.mode.(type) {
	case *modes.SwitchMode:
		if len(sessions) > 0 {
			m.mode = modes.NewRenameMode("")
		}
	case *modes.RenameMode:
		m.mode = modes.NewCreateMode(m.dirs)
	case *modes.CreateMode:
		m.mode = modes.NewSwitchMode(sessions)
	default:
		m.mode = modes.NewSwitchMode(sessions)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			getDefaultDirs returns a default set of directories.
//
//	@Return			[]string	Default directory paths
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) getDefaultDirs() []string {
	return []string{
		"~/projects",
		"~/code",
		"~/work",
		"~/.config",
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderHelpOverlay renders the help dialog as an overlay.
//
//	@Return			string	Help dialog overlaid on dimmed background
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderHelpOverlay() string {
	dim := styles.CurrentTheme.DimmedBackground.
		Width(m.totalContentWidth()).
		Height(styles.CurrentTheme.ContainerHeight).
		Render(strings.Repeat("\n", styles.CurrentTheme.ContainerHeight))
	help := RenderHelpDialog(m.totalContentWidth())
	dimmed := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dim)
	overlayed := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, help)
	return dimmed + "\n" + overlayed
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderHeader renders the application header.
//
//	@Return			string	Header with title and mode indicator
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderHeader() string {
	title := m.renderTitle()
	mode := m.renderModeIndicator()
	row := lipgloss.JoinHorizontal(lipgloss.Top, title, mode)
	return styles.CurrentTheme.HeaderStyle.Width(m.totalContentWidth()).Render(row)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderTitle renders the application title.
//
//	@Return			string	Styled "TSM" title
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderTitle() string {
	style := lipgloss.NewStyle().Bold(true).Foreground(styles.CurrentTheme.AccentColor)
	left := lipgloss.NewStyle().Width(styles.CurrentTheme.LeftPanelWidth).Render(style.Render("󱎫 TSM"))
	return left
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderModeIndicator renders the current mode indicator.
//
//	@Return			string	Styled mode name with icon
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderModeIndicator() string {
	modeText := m.mode.GetIcon() + " " + m.modeLabel()
	right := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Foreground(styles.CurrentTheme.HighlightColor).
		Bold(true).
		Width(styles.CurrentTheme.RightPanelWidth).
		Render(modeText)
	return right
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderBody renders the main content body.
//
//	@Return			string	Current mode's view content
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderBody() string {
	return styles.CurrentTheme.ListStyle.Render(m.mode.View())
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			renderFooter renders the application footer.
//
//	@Return			string	Footer with help text from current mode
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) renderFooter() string {
	text := m.mode.GetFooterText()
	styledText := lipgloss.NewStyle().Foreground(styles.CurrentTheme.SecondaryColor).Render(text)
	return styles.CurrentTheme.FooterStyle.Render(
		lipgloss.PlaceHorizontal(m.totalContentWidth(), lipgloss.Center, styledText),
	)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			remainingHeight calculates remaining vertical space for padding.
//
//	@Param			contentHeight	int		Current content height
//
//	@Return			int		Remaining height (0 if negative)
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) remainingHeight(contentHeight int) int {
	r := styles.CurrentTheme.ContainerHeight - contentHeight
	if r < 0 {
		return 0
	}
	return r
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			totalContentWidth calculates the total width of content area.
//
//	@Return			int		Total width (left panel + right panel + padding)
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) totalContentWidth() int {
	return styles.CurrentTheme.LeftPanelWidth + styles.CurrentTheme.RightPanelWidth + 2
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			modeLabel returns a cleaned-up mode label for display.
//
//	@Return			string	Formatted mode name
//
/////////////////////////////////////////////////////////////////////////////////////////////
func (m *manager) modeLabel() string {
	if m.mode == nil {
		sessions, _ := tmux.ListSessions()
		m.mode = modes.NewSwitchMode(sessions)
		return "SWITCH"
	}
	name := m.mode.ModeName()
	name = strings.TrimSuffix(name, " MODE")
	name = strings.TrimPrefix(name, "RENAME: ")
	name = strings.TrimSuffix(name, " - SELECT SESSION")
	return name
}
