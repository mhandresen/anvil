package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCandidatesOrder(t *testing.T) {
	env := map[string]string{
		"ANVIL_HOST":     "unix:///tmp/anvil.sock",
		"DOCKER_HOST":    "unix:///tmp/docker.sock",
		"CONTAINER_HOST": "unix:///tmp/podman.sock",
	}
	c := Candidates(1000, env)
	if len(c) == 0 || c[0] != "unix:///tmp/anvil.sock" {
		t.Fatalf("første kandidat %#v", c)
	}
	want := "unix:///run/user/1000/podman/podman.sock"
	found := false
	for _, s := range c {
		if s == want {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("mangler %s i %#v", want, c)
	}
}

func TestDetectNoSockets(t *testing.T) {
	_, err := DetectFrom(context.Background(), []string{
		"unix:///no/such/anvil.sock",
		"unix:///no/such/docker.sock",
		"unix:///no/such/podman.sock",
	})
	if err != ErrNoEngine {
		t.Fatalf("err=%v", err)
	}
}

func TestDetectFindsOverriddenSocket(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "podman.sock")
	if err := os.WriteFile(path, []byte("fake"), 0o600); err != nil {
		t.Fatal(err)
	}

	d, err := DetectFrom(context.Background(), []string{"unix://" + path})
	if err != nil {
		t.Fatal(err)
	}
	if d.Kind != KindPodman {
		t.Fatalf("kind=%s", d.Kind)
	}
	if d.Mode != ModeSystemPodman && d.Mode != ModeRootlessPodman {
		t.Fatalf("mode=%s", d.Mode)
	}
	if d.Socket != "unix://"+path {
		t.Fatalf("socket=%s", d.Socket)
	}
}