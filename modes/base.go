package modes

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ModeStrategy interface {
	Update(msg tea.Msg) (ModeStrategy, tea.Cmd)
	View() string
	ModeName() string
	Reset()
	GetCurrentSession() string
	GetIcon() string
	GetFooterText() string
}
