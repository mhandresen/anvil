package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := loadPaths()
		if err != nil {
			return err
		}
		servers, err := config.LoadServers(paths.Servers)
		if err != nil {
			return err
		}

		states := map[string]string{}
		if cli, _, err := openEngine(cmd.Context()); err == nil {
			defer cli.Close()
			for id := range servers.Servers {
				found, ferr := findNamed(cmd.Context(), cli, "anvil-"+id)
				if ferr != nil || found.ID == "" {
					states[id] = "-"
					continue
				}
				ins, ierr := cli.Inspect(cmd.Context(), found.ID)
				if ierr != nil || ins.State == "" {
					states[id] = found.State
					continue
				}
				states[id] = ins.State
			}
		}

		type row struct {
			ID     string `json:"id"`
			Game   string `json:"game"`
			Public bool   `json:"public"`
			State  string `json:"state"`
		}
		rows := make([]row, 0, len(servers.Servers))
		for id, s := range servers.Servers {
			st := states[id]
			if st == "" {
				st = "-"
			}
			rows = append(rows, row{ID: id, Game: s.Game, Public: s.Public, State: st})
		}

		if flagJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(rows)
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tGAME\tPUBLIC\tSTATE")
		for _, r := range rows {
			fmt.Fprintf(w, "%s\t%s\t%v\t%s\n", r.ID, r.Game, r.Public, r.State)
		}
		return w.Flush()
	},
}