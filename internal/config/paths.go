package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	ConfigDir    string
	DataDir      string
	StateDir     string
	AppConfig    string
	Servers      string
	Allocations  string
	PlayitSecret string
	PlayitHome   string
	Backups      string
	Runtime      string
}

func ResolvePaths(home string, env map[string]string) Paths {
	configDir := resolveDir(env["ANVIL_CONFIG_DIR"], env["XDG_CONFIG_HOME"], filepath.Join(home, ".config"), "anvil")
	dataDir := resolveDir(env["ANVIL_DATA_DIR"], env["XDG_DATA_HOME"], filepath.Join(home, ".local", "share"), "anvil")
	stateDir := resolveDir("", env["XDG_STATE_HOME"], filepath.Join(home, ".local", "state"), "anvil")

	return Paths{
		ConfigDir:    configDir,
		DataDir:      dataDir,
		StateDir:     stateDir,
		AppConfig:    filepath.Join(configDir, "config.toml"),
		Servers:      filepath.Join(configDir, "servers.toml"),
		Allocations:  filepath.Join(dataDir, "allocations.toml"),
		PlayitSecret: filepath.Join(configDir, "playit.secret"),
		PlayitHome:   filepath.Join(stateDir, "playit"),
		Backups:      filepath.Join(dataDir, "backups"),
		Runtime:      filepath.Join(stateDir, "runtime.json"),
	}
}

func resolveDir(override, xdg, fallbackParent, app string) string {
	if override != "" {
		return override
	}
	if xdg != "" {
		return filepath.Join(xdg, app)
	}
	return filepath.Join(fallbackParent, app)
}

func HomeDir() (string, error) {
	if h := os.Getenv("HOME"); h != "" {
		return h, nil
	}
	return os.UserHomeDir()
}