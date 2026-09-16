package config

type App struct {
	UI      UIConfig      `toml:"ui"`
	Engine  EngineConfig  `toml:"engine"`
	Network NetworkConfig `toml:"network"`
	Data    DataConfig    `toml:"data"`
	Tunnel  TunnelConfig  `toml:"tunnel"`
	Backup  BackupConfig  `toml:"backup"`
}

type UIConfig struct {
	Bind        string `toml:"bind"`
	OpenBrowser bool   `toml:"open_browser"`
}

type EngineConfig struct {
	Socket string `toml:"socket"`
}

type NetworkConfig struct {
	Subnet string `toml:"subnet"`
}

type DataConfig struct {
	Dir string `toml:"dir"`
}

type TunnelConfig struct {
	Provider         string      `toml:"provider"`
	PublicByDefault  bool        `toml:"public_by_default"`
	StopIdleAgent    bool        `toml:"stop_idle_agent"`
	Playit           PlayitConfig `toml:"playit"`
}

type PlayitConfig struct {
	Image  string `toml:"image"`
	Region string `toml:"region"`
}

type BackupConfig struct {
	Keep int `toml:"keep"`
}

func DefaultApp() App {
	return App{
		UI: UIConfig{
			Bind:        "127.0.0.1:8080",
			OpenBrowser: true,
		},
		Network: NetworkConfig{
			Subnet: "10.57.21.0/24",
		},
		Tunnel: TunnelConfig{
			Provider:      "playit",
			StopIdleAgent: true,
			Playit: PlayitConfig{
				Image: "ghcr.io/playit-cloud/playit-agent:latest",
			},
		},
		Backup: BackupConfig{
			Keep: 3,
		},
	}
}