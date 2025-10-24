package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jkeresman01/tsm/config"
	modes "github.com/jkeresman01/tsm/modes"
	styles "github.com/jkeresman01/tsm/styles"
	"github.com/jkeresman01/tsm/tmux"
	"github.com/jkeresman01/tsm/utils"
)

type manager struct {
	width    int
	height   int
	showHelp bool
	mode     modes.ModeStrategy
	dirs     []string
}

func NewTsmManager(cfg config.Config) tea.Model {
	sessions, _ := tmux.ListSessions()
	if len(sessions) == 0 {
		sessions = []string{}
	}

	// Load project directories from config
	dirs := utils.GetProjectDirs(cfg.SearchPaths, cfg.MaxDepth)

	return &manager{
		mode: modes.NewSwitchMode(sessions),
		dirs: dirs,
	}
}

func (m *manager) Init() tea.Cmd { return nil }

func (m *manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.WindowSizeMsg:
		m.applyWindowSize(t)
		return m, nil
	case tea.KeyMsg:
		// Handle global keys first
		if cmd := m.handleGlobalKey(t); cmd != nil {
			return m, cmd
		}
	}

	// Let the current mode handle the message
	newMode, cmd := m.mode.Update(msg)
	m.mode = newMode

	return m, cmd
}

func (m *manager) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.showHelp {
		return m.renderHelpOverlay()
	}

	header := m.renderHeader()
	body := m.renderBody()
	footer := m.renderFooter()

	contentHeight := lipgloss.Height(header) + lipgloss.Height(body) + lipgloss.Height(footer)
	padding := strings.Repeat("\n", m.remainingHeight(contentHeight))

	layout := lipgloss.JoinVertical(lipgloss.Top, header, body, padding, footer)
	withOuter := styles.CurrentTheme.OuterStyle.Render(layout)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, withOuter)
}

func (m *manager) applyWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
}

func (m *manager) handleGlobalKey(k tea.KeyMsg) tea.Cmd {
	switch k.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "?":
		m.showHelp = !m.showHelp
		return nil
	case "tab":
		m.cycleMode()
		return nil
	case "ctrl+n":
		m.mode = modes.NewCreateMode(m.dirs)
		// If dirs is empty, try to load some default directories
		if len(m.dirs) == 0 {
			m.dirs = m.getDefaultDirs()
			m.mode = modes.NewCreateMode(m.dirs)
		}
		return nil
	case "ctrl+r":
		// Always start rename mode in selection phase (empty string)
		m.mode = modes.NewRenameMode("")
		return nil
	case "ctrl+s":
		sessions, _ := tmux.ListSessions()
		m.mode = modes.NewSwitchMode(sessions)
		return nil
	}
	return nil
}

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

func (m *manager) getDefaultDirs() []string {
	return []string{
		"~/projects",
		"~/code",
		"~/work",
		"~/.config",
	}
}

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

func (m *manager) renderHeader() string {
	// Title with icon
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.CurrentTheme.AccentColor)

	title := titleStyle.Render("󱎫 TSM")

	// Mode indicator with icon
	modeIcon := m.getModeIcon()
	modeText := modeIcon + " " + m.modeLabel()

	mode := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Foreground(styles.CurrentTheme.HighlightColor).
		Bold(true).
		Render(modeText)

	left := lipgloss.NewStyle().
		Width(styles.CurrentTheme.LeftPanelWidth).
		Render(title)

	right := lipgloss.NewStyle().
		Width(styles.CurrentTheme.RightPanelWidth).
		Render(mode)

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return styles.CurrentTheme.HeaderStyle.
		Width(m.totalContentWidth()).
		Render(row)
}

func (m *manager) renderBody() string {
	return styles.CurrentTheme.ListStyle.Render(m.mode.View())
}

func (m *manager) renderFooter() string {
	text := m.getFooterText()

	styledText := lipgloss.NewStyle().
		Foreground(styles.CurrentTheme.SecondaryColor).
		Render(text)

	return styles.CurrentTheme.FooterStyle.Render(
		lipgloss.PlaceHorizontal(m.totalContentWidth(), lipgloss.Center, styledText),
	)
}

func (m *manager) getFooterText() string {
	switch m.mode.(type) {
	case *modes.SwitchMode:
		return "↑↓ navigate • ↵ switch • ⇥ cycle • ^N new • ^R rename • ? help • q quit"
	case *modes.RenameMode:
		return "↑↓ navigate • ↵ select/confirm • ⎋ cancel • ⇥ cycle • ? help • q quit"
	case *modes.CreateMode:
		return "↑↓ navigate • ↵ create • ⇥ cycle • ? help • q quit"
	default:
		return "? help • q quit"
	}
}

func (m *manager) getModeIcon() string {
	switch m.mode.(type) {
	case *modes.SwitchMode:
		return "󰆧" // Switch icon
	case *modes.RenameMode:
		return "󰑕" // Edit icon
	case *modes.CreateMode:
		return "󰐕" // Add icon
	default:
		return "󰍉" // Terminal icon
	}
}

func (m *manager) remainingHeight(contentHeight int) int {
	r := styles.CurrentTheme.ContainerHeight - contentHeight
	if r < 0 {
		return 0
	}
	return r
}

func (m *manager) totalContentWidth() int {
	return styles.CurrentTheme.LeftPanelWidth + styles.CurrentTheme.RightPanelWidth + 2
}

func (m *manager) modeLabel() string {
	if m.mode == nil {
		sessions, _ := tmux.ListSessions()
		m.mode = modes.NewSwitchMode(sessions)
		return "SWITCH"
	}

	// Simplify mode names
	name := m.mode.ModeName()
	name = strings.TrimSuffix(name, " MODE")
	name = strings.TrimPrefix(name, "RENAME: ")
	name = strings.TrimSuffix(name, " - SELECT SESSION")

	return name
}
