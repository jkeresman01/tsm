package tmux

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

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

func RenameSession(oldName, newName string) error {
	cmd := exec.Command("tmux", "rename-session", "-t", oldName, newName)
	return cmd.Run()
}

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

func CreateSession(name, path string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name, "-c", path)
	return cmd.Run()
}
