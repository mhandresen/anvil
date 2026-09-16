package config

import (
	"fmt"
	"regexp"
)

var serverIDRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,46}[a-z0-9])$`)

type ServersFile struct {
	Servers map[string]Server `toml:"servers"`
}

type Server struct {
	Game          string            `toml:"game"`
	Image         string            `toml:"image"`
	Public        bool              `toml:"public"`
	Memory        string            `toml:"memory"`
	CPUs          float64           `toml:"cpus"`
	PublishToHost *bool             `toml:"publish_to_host"`
	Env           map[string]string `toml:"env"`
	Ports         map[string]int    `toml:"ports"`
}

func (s ServersFile) Server(id string) (Server, bool) {
	if s.Servers == nil {
		return Server{}, false
	}
	v, ok := s.Servers[id]
	return v, ok
}

func ValidateServerID(id string) error {
	if !serverIDRe.MatchString(id) {
		return fmt.Errorf("invalid server-id %q", id)
	}
	return nil
}

func ValidateGame(game string) error {
	switch game {
	case "minecraft", "valheim", "palworld":
		return nil
	default:
		return fmt.Errorf("unknown game %q", game)
	}
}