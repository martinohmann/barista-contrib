package notify

import "os/exec"

// Send sends a new notification with title and content via notify-send.
func Send(title, content string) error {
	return exec.Command(
		"notify-send",
		"--expire-time", "10000",
		"--icon", "none",
		"--app-name", title,
		" ", content,
	).Run()
}
