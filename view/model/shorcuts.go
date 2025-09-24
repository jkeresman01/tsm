package model

type Shortcut struct {
	Key  string
	Desc string
}

var Shortcuts = []Shortcut{
	{"↑ / k", "Move up"},
	{"↓ / j", "Move down"},
	{"Enter", "Select / Confirm"},
	{"Tab", "Cycle mode"},
	{"Ctrl+N", "Create new session"},
	{"Ctrl+R", "Rename selected session"},
	{"Ctrl+S", "Switch to selected session"},
	{"q / Ctrl+C", "Quit"},
	{"?", "Toggle help"},
}
