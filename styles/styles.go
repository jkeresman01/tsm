package styles

import "github.com/charmbracelet/lipgloss"

import "strings"

const (
	LeftPanelWidth  = 40
	RightPanelWidth = 60
)

type Theme struct {
	BorderColor    lipgloss.Color
	SecondaryColor lipgloss.Color

	LeftPanelWidth  int
	RightPanelWidth int

	ContainerHeight int

	DimmedBackground lipgloss.Style

	OuterStyle   lipgloss.Style
	HeaderStyle  lipgloss.Style
	FooterStyle  lipgloss.Style
	ListStyle    lipgloss.Style
	PreviewStyle lipgloss.Style
}

var CurrentTheme Theme

func InitTheme(mode string) {
	CurrentTheme = pickTheme(normalizeMode(mode))
}

func normalizeMode(mode string) string {
	m := strings.TrimSpace(strings.ToLower(mode))
	if m == "" {
		return "dark"
	}
	return m
}

func pickTheme(mode string) Theme {
	switch mode {
	case "light":
		return LightTheme()
	case "dark":
		return DarkTheme()
	default:
		return DarkTheme()
	}
}
