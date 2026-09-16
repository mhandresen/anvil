package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const version = "dev"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version number of Anvil",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("anvil %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
	},
}