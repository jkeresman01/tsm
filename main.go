package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jkeresman01/tsm/logger_factory"
	"github.com/jkeresman01/tsm/styles"
	"github.com/jkeresman01/tsm/view"
)

func main() {
	styles.InitTheme("dark")

	log := logger_factory.GetLogger("tsm.log")

	p := tea.NewProgram(view.NewTsmMManager(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal("TSM exited with error:", err)
	}
}
