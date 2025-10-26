package styles

import "github.com/charmbracelet/lipgloss"

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			themeSizes defines dimensional properties for a theme.
//
/////////////////////////////////////////////////////////////////////////////////////////////
type themeSizes struct {
	leftPanel       int
	rightPanel      int
	containerHeight int
	paddingX        int
	paddingY        int
	marginX         int
	marginY         int
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			themePalette defines color properties for a theme.
//
/////////////////////////////////////////////////////////////////////////////////////////////
type themePalette struct {
	border    lipgloss.Color
	secondary lipgloss.Color
	dimBG     lipgloss.Color
	accent    lipgloss.Color
	highlight lipgloss.Color
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			baseBorderStyle creates a style with rounded borders.
//
//	@Param			border	lipgloss.Color	Border color
//
//	@Return			lipgloss.Style	Styled border
//
/////////////////////////////////////////////////////////////////////////////////////////////
func baseBorderStyle(border lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(border)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			outerStyle creates the style for the outer container.
//
//	@Param			p	themePalette	Color palette
//	@Param			s	themeSizes		Size dimensions
//
//	@Return			lipgloss.Style	Outer container style
//
/////////////////////////////////////////////////////////////////////////////////////////////
func outerStyle(p themePalette, s themeSizes) lipgloss.Style {
	return baseBorderStyle(p.border).
		BorderStyle(lipgloss.RoundedBorder()).
		Height(s.containerHeight).
		Padding(s.paddingY, s.paddingX).
		Margin(s.marginY, s.marginX)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			headerStyle creates the style for the header.
//
//	@Param			p	themePalette	Color palette
//
//	@Return			lipgloss.Style	Header style
//
/////////////////////////////////////////////////////////////////////////////////////////////
func headerStyle(p themePalette) lipgloss.Style {
	return baseBorderStyle(p.border).
		BorderBottom(true).
		Padding(0, 1).
		MarginBottom(1)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			footerStyle creates the style for the footer.
//
//	@Param			p	themePalette	Color palette
//
//	@Return			lipgloss.Style	Footer style
//
/////////////////////////////////////////////////////////////////////////////////////////////
func footerStyle(p themePalette) lipgloss.Style {
	return baseBorderStyle(p.border).
		BorderTop(true).
		Foreground(p.secondary).
		Padding(0, 1).
		MarginTop(1)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			dimStyle creates the style for dimmed backgrounds.
//
//	@Param			p	themePalette	Color palette
//
//	@Return			lipgloss.Style	Dimmed background style
//
/////////////////////////////////////////////////////////////////////////////////////////////
func dimStyle(p themePalette) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(p.dimBG)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			listStyle creates the style for the list panel.
//
//	@Param			s	themeSizes	Size dimensions
//
//	@Return			lipgloss.Style	List panel style
//
/////////////////////////////////////////////////////////////////////////////////////////////
func listStyle(s themeSizes) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(s.leftPanel).
		Padding(0, 1)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			previewStyle creates the style for the preview panel.
//
//	@Param			s	themeSizes	Size dimensions
//
//	@Return			lipgloss.Style	Preview panel style
//
/////////////////////////////////////////////////////////////////////////////////////////////
func previewStyle(s themeSizes) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(s.rightPanel).
		Padding(0, 1)
}
