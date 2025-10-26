package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jkeresman01/tsm/config"
	"github.com/jkeresman01/tsm/logger_factory"
	"github.com/jkeresman01/tsm/styles"
	"github.com/jkeresman01/tsm/view"
)

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			main initializes and starts the TSM application.
//
//	@Description	Loads configuration from ~/.config/tsm/config.json
//	@Description	Creates default config if none exists
//	@Description	Initializes UI theme based on configuration
//	@Description	Sets up logging to tsm.log
//	@Description	Starts the Bubble Tea TUI program
//
/////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	cfg, err := config.Load()
	if err != nil {
		config.SaveDefault()
		cfg = config.DefaultConfig()
	}

	styles.InitTheme(cfg.Theme)

	log := logger_factory.GetLogger("tsm.log")

	p := tea.NewProgram(view.NewTsmManager(cfg), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal("TSM exited with error:", err)
	}
}
