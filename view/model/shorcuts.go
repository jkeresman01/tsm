package model

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			Shortcut represents a keyboard shortcut with its description.
//
/////////////////////////////////////////////////////////////////////////////////////////////
type Shortcut struct {
	Key  string
	Desc string
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief Shortcuts contains all keyboard shortcuts for the application.
//
//	@Description	Displayed in the help dialog when user presses '?'
//
/////////////////////////////////////////////////////////////////////////////////////////////
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
