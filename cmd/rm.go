package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/spf13/cobra"
)

var rmWipe bool

var rmCmd = &cobra.Command{
	Use:   "rm <id>",
	Short: "Remove a server from config and stop its container",
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
		if _, ok := servers.Server(id); !ok {
			return fmt.Errorf("unknown server %s", id)
		}

		if cli, _, err := openEngine(cmd.Context()); err == nil {
			defer cli.Close()
			found, ferr := findNamed(cmd.Context(), cli, "anvil-"+id)
			if ferr == nil && found.ID != "" {
				_ = cli.Stop(cmd.Context(), found.ID, 30*time.Second)
				if rerr := cli.Remove(cmd.Context(), found.ID, false); rerr != nil {
					return rerr
				}
			}
		}

		delete(servers.Servers, id)
		if err := config.SaveServers(paths.Servers, servers); err != nil {
			return err
		}

		allocs, err := config.LoadAllocations(paths.Allocations)
		if err == nil {
			delete(allocs.Servers, id)
			_ = config.SaveAllocations(paths.Allocations, allocs)
		}

		if rmWipe {
			dir := filepath.Join(paths.DataDir, "data", id)
			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "deleted %s\n", id)
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVar(&rmWipe, "wipe", false, "also delete the data directory")
}