package engine

import (
	"bytes"
	"fmt"
	"os/exec"
)

func EnableLinger() error {
	cmd := exec.Command("loginctl", "enable-linger")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("loginctl enable-linger: %w (%s)", err, bytes.TrimSpace(out))
	}
	return nil
}
