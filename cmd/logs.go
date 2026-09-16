package cmd

import (
	"fmt"
	"os"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/spf13/cobra"
)

var logsFollow bool

var logsCmd = &cobra.Command{
	Use:   "logs <id>",
	Short: "Show server logs",
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
		rc, err := cli.Logs(cmd.Context(), found.ID, logsFollow)
		if err != nil {
			return err
		}
		defer rc.Close()
		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, rc)
		return err
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "follow log output")
}