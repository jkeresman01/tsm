package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jkeresman01/tsm/styles"
	"github.com/jkeresman01/tsm/view/model"
)

const helpKeyColWidth = 14

func RenderHelpDialog(width int) string {
	content := helpContent()
	return styles.HelpBoxStyle.Width(width).Render(content)
}

func helpContent() string {
	var b strings.Builder
	b.WriteString(helpTitle())
	b.WriteString("\n\n")
	b.WriteString(helpLines())
	return b.String()
}

func helpTitle() string {
	return styles.HelpTitleStyle.Render("Keyboard Shortcuts")
}

func helpLines() string {
	var b strings.Builder
	for _, sc := range model.Shortcuts {
		b.WriteString(helpLine(sc))
		b.WriteByte('\n')
	}
	return b.String()
}

func helpLine(sc model.Shortcut) string {
	key := styles.HelpKeyStyle.Width(helpKeyColWidth).Render(sc.Key)
	desc := styles.HelpDescStyle.Render(sc.Desc)
	return lipgloss.JoinHorizontal(lipgloss.Top, key, desc)
}
