package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
)

type Allocations struct {
	HostID  string                  `toml:"host_id"`
	Network AllocNetwork            `toml:"network"`
	Servers map[string]ServerAlloc  `toml:"servers"`
	Agent   AgentAlloc              `toml:"agent"`
}

type AllocNetwork struct {
	Subnet string `toml:"subnet"`
}

type ServerAlloc struct {
	IP               string `toml:"ip"`
	DesiredState     string `toml:"desired_state"`
	PlayitTunnelID   string `toml:"playit_tunnel_id"`
	PlayitTunnelName string `toml:"playit_tunnel_name"`
}

type AgentAlloc struct {
	IP string `toml:"ip"`
}

func EmptyAllocations(subnet string) Allocations {
	return Allocations{
		Network: AllocNetwork{Subnet: subnet},
		Servers: map[string]ServerAlloc{},
		Agent:   AgentAlloc{IP: ""},
	}
}

func NewHostID() (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func ParseSubnet(cidr string) (*net.IPNet, error) {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet %q: %w", cidr, err)
	}
	return n, nil
}