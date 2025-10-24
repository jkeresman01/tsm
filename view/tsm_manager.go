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
		if cmd := m.handleGlobalKey(t); cmd != nil {
			return m, cmd
		}
	}

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
		if len(m.dirs) == 0 {
			m.dirs = m.getDefaultDirs()
			m.mode = modes.NewCreateMode(m.dirs)
		}
		return nil
	case "ctrl+r":
		var sessionToRename string
		if switchMode, ok := m.mode.(*modes.SwitchMode); ok {
			sessionToRename = switchMode.GetCurrentSession()
		}

		if sessionToRename == "" {
			sessions, _ := tmux.ListSessions()
			if len(sessions) > 0 {
				sessionToRename = sessions[0]
			}
		}

		if sessionToRename != "" {
			m.mode = modes.NewRenameMode(sessionToRename)
		}
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
			m.mode = modes.NewRenameMode(sessions[0])
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
		"~/git",
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
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.CurrentTheme.BorderColor).
		Render("TSM")

	left := lipgloss.NewStyle().
		Width(styles.CurrentTheme.LeftPanelWidth).
		Render(title)

	mode := lipgloss.NewStyle().
		Width(styles.CurrentTheme.RightPanelWidth).
		Align(lipgloss.Right).
		Foreground(styles.CurrentTheme.SecondaryColor).
		Render(m.modeLabel())

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, mode)

	return styles.CurrentTheme.HeaderStyle.
		Width(m.totalContentWidth()).
		Render(row)
}

func (m *manager) renderBody() string {
	return styles.CurrentTheme.ListStyle.Render(m.mode.View())
}

func (m *manager) renderFooter() string {
	text := m.getFooterText()
	return styles.CurrentTheme.FooterStyle.Render(
		lipgloss.PlaceHorizontal(m.totalContentWidth(), lipgloss.Center, text),
	)
}

func (m *manager) getFooterText() string {
	switch m.mode.(type) {
	case *modes.SwitchMode:
		return "↑↓ to navigate • enter to switch • tab to cycle mode • ? for help • q to quit"
	case *modes.RenameMode:
		return "type new name • enter to confirm • tab to cycle mode • q to quit"
	case *modes.CreateMode:
		return "↑↓ to navigate • enter to create session • tab to cycle mode • q to quit"
	default:
		return "? for help • q to quit"
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
		return "SWITCH MODE"
	}
	return m.mode.ModeName()
}
