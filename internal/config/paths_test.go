package config

import (
	"path/filepath"
	"testing"
)

func TestResolvePathsDefaults(t *testing.T) {
	p := ResolvePaths("/home/dev", map[string]string{})
	if p.ConfigDir != "/home/dev/.config/anvil" {
		t.Fatalf("config=%s", p.ConfigDir)
	}
	if p.Allocations != "/home/dev/.local/share/anvil/allocations.toml" {
		t.Fatalf("alloc=%s", p.Allocations)
	}
	if p.Servers != "/home/dev/.config/anvil/servers.toml" {
		t.Fatalf("servers=%s", p.Servers)
	}
}

func TestResolvePathsOverrides(t *testing.T) {
	p := ResolvePaths("/home/dev", map[string]string{
		"ANVIL_CONFIG_DIR": "/tmp/cfg",
		"ANVIL_DATA_DIR":   "/tmp/data",
	})
	if p.ConfigDir != "/tmp/cfg" {
		t.Fatalf("config=%s", p.ConfigDir)
	}
	if p.Allocations != filepath.Join("/tmp/data", "allocations.toml") {
		t.Fatalf("alloc=%s", p.Allocations)
	}
}

func TestAllocationsNotUnderConfig(t *testing.T) {
	p := ResolvePaths("/home/dev", map[string]string{})
	if filepath.Dir(p.Allocations) == p.ConfigDir {
		t.Fatal("allocations skal ikke ligge i config-dir")
	}
}