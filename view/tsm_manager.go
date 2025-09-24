package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jkeresman01/tsm/view/model"

	tea "github.com/charmbracelet/bubbletea"
	modes "github.com/jkeresman01/tsm/modes"
	styles "github.com/jkeresman01/tsm/styles"
)

type manager struct {
	width        int
	height       int
	showHelp     bool
	sessionModel model.SessionModel
}

func NewTsmMManager() tea.Model {
	return &manager{
		sessionModel: model.NewSessionModel(),
	}
}

func (m *manager) Init() tea.Cmd { return nil }

func (m *manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.WindowSizeMsg:
		m.applyWindowSize(t)
		return m, nil
	case tea.KeyMsg:
		m.handleKey(t)
	}

	updated, cmd := m.sessionModel.Update(msg)
	m.sessionModel = updated
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

func (m *manager) handleKey(k tea.KeyMsg) {
	switch k.String() {
	case "tab":
		m.sessionModel.NextMode()
	case "ctrl+n":
		m.sessionModel.SetMode(modes.NewCreateMode(m.sessionModel.Sessions()))
	case "ctrl:r", "ctrl+r":
		m.sessionModel.SetMode(modes.NewRenameMode(m.sessionModel.CurrentSession()))
	case "ctrl:s", "ctrl+s":
		m.sessionModel.SetMode(modes.NewSwitchMode(m.sessionModel.Sessions()))
	case "?":
		m.showHelp = !m.showHelp
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

	search := m.sessionModel.Input().View()

	left := lipgloss.NewStyle().
		Width(styles.CurrentTheme.LeftPanelWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, title, " 🔍 ", search))

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
	return styles.CurrentTheme.ListStyle.Render(m.sessionModel.View())
}

func (m *manager) renderFooter() string {
	text := "↑↓ to navigate • enter to attach • q to quit"
	return styles.CurrentTheme.FooterStyle.Render(
		lipgloss.PlaceHorizontal(m.totalContentWidth(), lipgloss.Center, text),
	)
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
	if m.sessionModel.Mode == nil {
		m.sessionModel.SetMode(modes.NewSwitchMode(m.sessionModel.Sessions()))
		return "SWITCH MODE"
	}
	return m.sessionModel.Mode.ModeName()
}
