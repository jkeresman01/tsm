package styles

import "github.com/charmbracelet/lipgloss"

import "strings"

const (
	LeftPanelWidth  = 45
	RightPanelWidth = 55
)

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			Theme defines the visual styling for the application.
//
//	@Description	Contains colors, dimensions, and lipgloss styles for UI components
//
/////////////////////////////////////////////////////////////////////////////////////////////
type Theme struct {
	BorderColor    lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor    lipgloss.Color
	HighlightColor lipgloss.Color

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

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			InitTheme initializes the theme based on the provided mode.
//
//	@Description	Loads "dark" or "light" theme and refreshes help styles
//
//	@Param			mode	string	Theme mode ("dark" or "light")
//
/////////////////////////////////////////////////////////////////////////////////////////////
func InitTheme(mode string) {
	CurrentTheme = pickTheme(normalizeMode(mode))
	RefreshHelpStyles()
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			normalizeMode normalizes the theme mode string.
//
//	@Description	Trims whitespace, converts to lowercase, defaults to "dark"
//
//	@Param			mode	string	Raw theme mode string
//
//	@Return			string	Normalized theme mode
//
/////////////////////////////////////////////////////////////////////////////////////////////
func normalizeMode(mode string) string {
	m := strings.TrimSpace(strings.ToLower(mode))
	if m == "" {
		return "dark"
	}
	return m
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			pickTheme selects the appropriate theme based on mode.
//
//	@Param			mode	string	Theme mode ("dark" or "light")
//
//	@Return			Theme	Selected theme (defaults to dark)
//
/////////////////////////////////////////////////////////////////////////////////////////////
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
