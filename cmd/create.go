package cmd

import (
	"fmt"
	"os"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/spf13/cobra"
)

var (
	createGame       string
	createPublic     bool
	createMemory     string
	createAcceptEULA bool
	createVersion    string
)

var createCmd = &cobra.Command{
	Use:   "create <id>",
	Short: "Create a server in config (does not start it)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if err := config.ValidateServerID(id); err != nil {
			return err
		}
		if err := config.ValidateGame(createGame); err != nil {
			return err
		}
		if createGame == "minecraft" && !createAcceptEULA {
			return fmt.Errorf("Minecraft requires --accept-eula")
		}

		paths, err := loadPaths()
		if err != nil {
			return err
		}
		servers, err := config.LoadServers(paths.Servers)
		if err != nil {
			return err
		}
		if _, exists := servers.Server(id); exists {
			return fmt.Errorf("server %s already exists", id)
		}
		if servers.Servers == nil {
			servers.Servers = map[string]config.Server{}
		}

		s := config.Server{
			Game:   createGame,
			Public: createPublic,
			Memory: createMemory,
			Env:    map[string]string{},
			Ports:  map[string]int{},
		}
		switch createGame {
		case "minecraft":
			s.Image = "docker.io/itzg/minecraft-server:latest"
			s.Env["EULA"] = "TRUE"
			s.Env["TYPE"] = "PAPER"
			if createVersion == "" {
				createVersion = "1.21.4"
			}
			s.Env["VERSION"] = createVersion
			s.Ports["game"] = 25565
		case "valheim":
			s.Image = "docker.io/lloesche/valheim-server:latest"
			s.Ports["game"] = 2456
		case "palworld":
			s.Image = "docker.io/jammsen/palworld-dedicated-server:latest"
			s.Ports["game"] = 8211
		}

		servers.Servers[id] = s
		if err := config.SaveServers(paths.Servers, servers); err != nil {
			return err
		}
		data := fmt.Sprintf("%s/data/%s", paths.DataDir, id)
		if err := os.MkdirAll(data, 0o700); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "created %s (%s)\n", id, createGame)
		fmt.Fprintf(cmd.OutOrStdout(), "config: %s\n", paths.Servers)
		fmt.Fprintf(cmd.OutOrStdout(), "data:   %s\n", data)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&createGame, "game", "", "minecraft|valheim|palworld")
	createCmd.Flags().BoolVar(&createPublic, "public", false, "public tunnel (later)")
	createCmd.Flags().StringVar(&createMemory, "memory", "4G", "memory limit")
	createCmd.Flags().BoolVar(&createAcceptEULA, "accept-eula", false, "accept the Minecraft EULA")
	createCmd.Flags().StringVar(&createVersion, "version", "", "Minecraft version (pinned in config)")
	_ = createCmd.MarkFlagRequired("game")
}