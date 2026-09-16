package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/mhandresen/anvil/internal/engine"
	"github.com/mhandresen/anvil/internal/tunnel"
	"github.com/spf13/cobra"
)

var tunnelSecret string

var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "Manage the Playit tunnel agent",
}

var tunnelLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Store a Playit agent secret and start the agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := strings.TrimSpace(tunnelSecret)
		if secret == "" {
			fmt.Fprintln(cmd.ErrOrStderr(), "Create a Docker agent at https://playit.gg/account/agents/new-docker")
			fmt.Fprint(cmd.ErrOrStderr(), "Paste SECRET_KEY: ")
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return err
			}
			secret = strings.TrimSpace(line)
		}
		if secret == "" {
			return fmt.Errorf("empty secret")
		}

		paths, err := loadPaths()
		if err != nil {
			return err
		}
		if err := tunnel.WriteSecret(paths.PlayitSecret, secret); err != nil {
			return err
		}

		cli, detected, err := openEngine(cmd.Context())
		if err != nil {
			return err
		}
		defer cli.Close()

		app, err := config.LoadApp(paths.AppConfig)
		if err != nil {
			return err
		}
		if _, err := cli.EnsureNetwork(cmd.Context(), engine.NetworkSpec{
			Name:   "anvil",
			Subnet: app.Network.Subnet,
		}); err != nil {
			return err
		}

		found, err := findNamed(cmd.Context(), cli, tunnel.AgentName)
		if err != nil {
			return err
		}
		if found.ID == "" {
			if !cli.ImageExists(cmd.Context(), tunnel.AgentImage) {
				fmt.Fprintf(cmd.OutOrStdout(), "pulling %s\n", tunnel.AgentImage)
				if err := cli.Pull(cmd.Context(), tunnel.AgentImage, cmd.OutOrStdout()); err != nil {
					return err
				}
			}
			id, err := cli.Create(cmd.Context(), engine.ContainerSpec{
				Name:    tunnel.AgentName,
				Image:   tunnel.AgentImage,
				Network: "anvil",
				Env:     map[string]string{"SECRET_KEY": secret},
				Labels: map[string]string{
					"io.anvil.managed": "true",
					"io.anvil.role":    "tunnel",
					"io.anvil.scope":   "host",
				},
			})
			if err != nil {
				return err
			}
			found.ID = id
			_ = detected
		}
		if err := cli.Start(cmd.Context(), found.ID); err != nil {
			return err
		}
		ins, err := cli.Inspect(cmd.Context(), found.ID)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "agent %s\n", ins.State)
		fmt.Fprintf(cmd.OutOrStdout(), "secret %s\n", paths.PlayitSecret)
		return nil
	},
}

var tunnelStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Playit agent status",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := loadPaths()
		if err != nil {
			return err
		}
		_, err = tunnel.ReadSecret(paths.PlayitSecret)
		hasSecret := err == nil
		fmt.Fprintf(cmd.OutOrStdout(), "secret:  %v\n", hasSecret)

		cli, _, err := openEngine(cmd.Context())
		if err != nil {
			return err
		}
		defer cli.Close()
		found, err := findNamed(cmd.Context(), cli, tunnel.AgentName)
		if err != nil {
			return err
		}
		if found.ID == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "agent:   missing")
			return nil
		}
		ins, err := cli.Inspect(cmd.Context(), found.ID)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "agent:   %s\n", ins.State)
		if ins.IP != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "ip:      %s\n", ins.IP)
		}
		return nil
	},
}

func init() {
	tunnelLoginCmd.Flags().StringVar(&tunnelSecret, "secret", "", "Playit agent SECRET_KEY")
	tunnelCmd.AddCommand(tunnelLoginCmd)
	tunnelCmd.AddCommand(tunnelStatusCmd)
}
