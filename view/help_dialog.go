package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jkeresman01/tsm/styles"
	"github.com/jkeresman01/tsm/view/model"
)

const helpKeyColWidth = 14

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			RenderHelpDialog renders the help dialog with keyboard shortcuts.
//
//		@Param			width	int		Width of the dialog
//
//		@Return			string	Rendered help dialog
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func RenderHelpDialog(width int) string {
	content := helpContent()
	return styles.HelpBoxStyle.Width(width).Render(content)
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			helpContent generates the complete help dialog content.
//
//		@Return			string	Help dialog content with title and shortcuts
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func helpContent() string {
	var b strings.Builder
	b.WriteString(helpTitle())
	b.WriteString("\n\n")
	b.WriteString(helpLines())
	return b.String()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			helpTitle renders the help dialog title.
//
//		@Return			string	Styled title "Keyboard Shortcuts"
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func helpTitle() string {
	return styles.HelpTitleStyle.Render("Keyboard Shortcuts")
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			helpLines renders all shortcut lines.
//
//		@Return			string	All shortcut lines formatted
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func helpLines() string {
	var b strings.Builder
	for _, sc := range model.Shortcuts {
		b.WriteString(helpLine(sc))
		b.WriteByte('\n')
	}
	return b.String()
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			helpLine renders a single shortcut line.
//
//		@Param			sc	model.Shortcut	Shortcut to render
//
//		@Return			string	Formatted shortcut line with key and description
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func helpLine(sc model.Shortcut) string {
	key := styles.HelpKeyStyle.Width(helpKeyColWidth).Render(sc.Key)
	desc := styles.HelpDescStyle.Render(sc.Desc)
	return lipgloss.JoinHorizontal(lipgloss.Top, key, desc)
}
