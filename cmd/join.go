package cmd

import (
	"fmt"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join <id>",
	Short: "Print the address players should use",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if err := config.ValidateServerID(id); err != nil {
			return err
		}
		paths, err := loadPaths()
		if err != nil {
			return err
		}
		servers, err := config.LoadServers(paths.Servers)
		if err != nil {
			return err
		}
		srv, ok := servers.Server(id)
		if !ok {
			return fmt.Errorf("unknown server %s", id)
		}

		port := srv.Ports["game"]
		if port == 0 {
			port = 25565
		}

		if srv.Public {
			fmt.Fprintf(cmd.OutOrStdout(), "public tunnel is not wired yet\n")
			fmt.Fprintf(cmd.OutOrStdout(), "lan: localhost:%d\n", port)
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "localhost:%d\n", port)
		return nil
	},
}