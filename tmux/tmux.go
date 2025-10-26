package tmux

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			ListSessions retrieves all active tmux sessions.
//
//	@Description	Executes 'tmux list-sessions' and parses session names
//
//	@Return			[]string	List of session names
//	@Return			error		Error if tmux command fails
//
/////////////////////////////////////////////////////////////////////////////////////////////
func ListSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	return lines, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			GetPreview captures the last 10 lines of a tmux session's active pane.
//
//	@Description	Executes 'tmux capture-pane' to get preview content
//
//	@Param			session	string	Session name to preview
//
//	@Return			string	Captured pane content
//	@Return			error	Error if tmux command fails
//
/////////////////////////////////////////////////////////////////////////////////////////////
func GetPreview(session string) (string, error) {
	cmd := exec.Command("tmux", "capture-pane", "-t", session, "-p", "-S", "-10")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			RenameSession renames an existing tmux session.
//
//	@Description	Executes 'tmux rename-session'
//
//	@Param			oldName	string	Current session name
//	@Param			newName	string	New session name
//
//	@Return			error	Error if tmux command fails
//
/////////////////////////////////////////////////////////////////////////////////////////////
func RenameSession(oldName, newName string) error {
	cmd := exec.Command("tmux", "rename-session", "-t", oldName, newName)
	return cmd.Run()
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			AttachSession attaches to or switches to a tmux session.
//
//	@Description	Uses 'switch-client' if already in tmux, otherwise 'attach-session'
//
//	@Param			name	string	Session name to attach
//
//	@Return			error	Error if tmux command fails
//
/////////////////////////////////////////////////////////////////////////////////////////////
func AttachSession(name string) error {
	if os.Getenv("TMUX") != "" {
		cmd := exec.Command("tmux", "switch-client", "-t", name)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmd := exec.Command("tmux", "attach-session", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			CreateSession creates a new detached tmux session.
//
//	@Description	Executes 'tmux new-session' with specified name and working directory
//
//	@Param			name	string	Session name
//	@Param			path	string	Working directory for the session
//
//	@Return			error	Error if tmux command fails
//
/////////////////////////////////////////////////////////////////////////////////////////////
func CreateSession(name, path string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name, "-c", path)
	return cmd.Run()
}
