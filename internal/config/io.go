package config

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

func LoadApp(path string) (App, error) {
	app := DefaultApp()
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return app, nil
		}
		return App{}, err
	}
	if err := toml.Unmarshal(b, &app); err != nil {
		return App{}, fmt.Errorf("config.toml: %w", err)
	}
	return app, nil
}

func SaveApp(path string, app App) error {
	return writeTOML(path, app)
}

func LoadServers(path string) (ServersFile, error) {
	var f ServersFile
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ServersFile{Servers: map[string]Server{}}, nil
		}
		return ServersFile{}, err
	}
	if err := toml.Unmarshal(b, &f); err != nil {
		return ServersFile{}, fmt.Errorf("servers.toml: %w", err)
	}
	if f.Servers == nil {
		f.Servers = map[string]Server{}
	}
	return f, nil
}

func SaveServers(path string, f ServersFile) error {
	if f.Servers == nil {
		f.Servers = map[string]Server{}
	}
	return writeTOML(path, f)
}

func LoadAllocations(path string) (Allocations, error) {
	var a Allocations
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return EmptyAllocations("10.57.21.0/24"), nil
		}
		return Allocations{}, err
	}
	if err := toml.Unmarshal(b, &a); err != nil {
		return Allocations{}, fmt.Errorf("allocations.toml: %w", err)
	}
	if a.Servers == nil {
		a.Servers = map[string]ServerAlloc{}
	}
	return a, nil
}

func SaveAllocations(path string, a Allocations) error {
	if a.Servers == nil {
		a.Servers = map[string]ServerAlloc{}
	}
	return writeTOML(path, a)
}

func writeTOML(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	b, err := toml.Marshal(v)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}