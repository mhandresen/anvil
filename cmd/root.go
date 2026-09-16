package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/mhandresen/anvil/internal/engine"
)

var (
	flagVerbose bool
	flagJSON    bool
)

var rootCmd = &cobra.Command{
	Use:           "anvil",
	Short:         "Container-agnostic game server manager for Linux",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(detectCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(joinCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		if errors.Is(err, engine.ErrNoEngine) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}