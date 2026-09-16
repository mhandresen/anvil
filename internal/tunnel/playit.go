package tunnel

import (
	"os"
	"strings"
)

const (
	AgentName  = "anvil-playit"
	AgentImage = "ghcr.io/playit-cloud/playit-agent:latest"
)

func ReadSecret(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func WriteSecret(path, secret string) error {
	if err := os.MkdirAll(dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.TrimSpace(secret)+"\n"), 0o600)
}

func dir(path string) string {
	i := strings.LastIndex(path, "/")
	if i <= 0 {
		return "."
	}
	return path[:i]
}
