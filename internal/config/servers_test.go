package config

import (
	"path/filepath"
	"testing"
)

func TestValidateServerID(t *testing.T) {
	ok := []string{"survival", "ab", "mc-1"}
	for _, id := range ok {
		if err := ValidateServerID(id); err != nil {
			t.Fatalf("%s: %v", id, err)
		}
	}
	bad := []string{"", "-x", "X", "a_", "Sur"}
	for _, id := range bad {
		if err := ValidateServerID(id); err == nil {
			t.Fatalf("forventet feil for %q", id)
		}
	}
}

func TestServersRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "servers.toml")

	in := ServersFile{Servers: map[string]Server{
		"survival": {
			Game:   "minecraft",
			Image:  "docker.io/itzg/minecraft-server:latest",
			Public: true,
			Memory: "4G",
			Env:    map[string]string{"EULA": "TRUE", "VERSION": "1.21.4"},
		},
	}}
	if err := SaveServers(path, in); err != nil {
		t.Fatal(err)
	}
	out, err := LoadServers(path)
	if err != nil {
		t.Fatal(err)
	}
	s, ok := out.Server("survival")
	if !ok || s.Game != "minecraft" || s.Env["VERSION"] != "1.21.4" {
		t.Fatalf("%#v", out)
	}
}

func TestLoadMissingFilesUseDefaults(t *testing.T) {
	dir := t.TempDir()
	app, err := LoadApp(filepath.Join(dir, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if app.Network.Subnet != "10.57.21.0/24" {
		t.Fatalf("subnet=%s", app.Network.Subnet)
	}
	a, err := LoadAllocations(filepath.Join(dir, "allocations.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if a.Servers == nil {
		t.Fatal("servers map nil")
	}
}