package styles

import "github.com/charmbracelet/lipgloss"

type themeSizes struct {
	leftPanel       int
	rightPanel      int
	containerHeight int
	paddingX        int
	paddingY        int
	marginX         int
	marginY         int
}

type themePalette struct {
	border    lipgloss.Color
	secondary lipgloss.Color
	dimBG     lipgloss.Color
	accent    lipgloss.Color
	highlight lipgloss.Color
}

func baseBorderStyle(border lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(border)
}

func outerStyle(p themePalette, s themeSizes) lipgloss.Style {
	return baseBorderStyle(p.border).
		BorderStyle(lipgloss.RoundedBorder()).
		Height(s.containerHeight).
		Padding(s.paddingY, s.paddingX).
		Margin(s.marginY, s.marginX)
}

func headerStyle(p themePalette) lipgloss.Style {
	return baseBorderStyle(p.border).
		BorderBottom(true).
		Padding(0, 1).
		MarginBottom(1)
}

func footerStyle(p themePalette) lipgloss.Style {
	return baseBorderStyle(p.border).
		BorderTop(true).
		Foreground(p.secondary).
		Padding(0, 1).
		MarginTop(1)
}

func dimStyle(p themePalette) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(p.dimBG)
}

func listStyle(s themeSizes) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(s.leftPanel).
		Padding(0, 1)
}

func previewStyle(s themeSizes) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(s.rightPanel).
		Padding(0, 1)
}
