package styles

import "github.com/charmbracelet/lipgloss"

var (
	HelpKeyStyle   = helpKey()
	HelpDescStyle  = helpDesc()
	HelpTitleStyle = helpTitle()
	HelpBoxStyle   = helpBox()
)

func RefreshHelpStyles() {
	HelpKeyStyle = helpKey()
	HelpDescStyle = helpDesc()
	HelpTitleStyle = helpTitle()
	HelpBoxStyle = helpBox()
}

func accent() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(CurrentTheme.BorderColor)
}

func secondary() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(CurrentTheme.SecondaryColor)
}

func bordered() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(CurrentTheme.BorderColor)
}

func helpKey() lipgloss.Style {
	return accent().Bold(true)
}

func helpDesc() lipgloss.Style {
	return secondary()
}

func helpTitle() lipgloss.Style {
	return accent().MarginBottom(1).Bold(true).Align(lipgloss.Center)
}

func helpBox() lipgloss.Style {
	totalWidth := CurrentTheme.LeftPanelWidth + CurrentTheme.RightPanelWidth
	return bordered().Padding(1, 3).Margin(1).Width(totalWidth)
}
