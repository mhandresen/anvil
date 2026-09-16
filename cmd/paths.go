package cmd

import "github.com/mhandresen/anvil/internal/config"

func loadPaths() (config.Paths, error) {
	home, err := config.HomeDir()
	if err != nil {
		return config.Paths{}, err
	}
	return config.ResolvePaths(home, envMap()), nil
}