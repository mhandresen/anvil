package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/mhandresen/anvil/internal/config"
)

var stopCmd = &cobra.Command{
	Use:   "stop <id>",
	Short: "Stop a server container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if err := config.ValidateServerID(id); err != nil {
			return err
		}
		cli, _, err := openEngine(cmd.Context())
		if err != nil {
			return err
		}
		defer cli.Close()

		found, err := findNamed(cmd.Context(), cli, "anvil-"+id)
		if err != nil {
			return err
		}
		if found.ID == "" {
			return fmt.Errorf("unknown server %s", id)
		}
		if err := cli.Stop(cmd.Context(), found.ID, 30*time.Second); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "stopped %s\n", id)
		return nil
	},
}