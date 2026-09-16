package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mhandresen/anvil/internal/config"
	"github.com/mhandresen/anvil/internal/engine"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <id>",
	Short: "Start a server container",
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

		cli, detected, err := openEngine(cmd.Context())
		if err != nil {
			return err
		}
		defer cli.Close()

		if detected.Linger != nil && !*detected.Linger {
			fmt.Fprintln(cmd.ErrOrStderr(), "warning: linger is disabled; the server will die on logout")
			fmt.Fprintln(cmd.ErrOrStderr(), "  loginctl enable-linger $USER")
		}

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

		if cli.ImageExists(cmd.Context(), srv.Image) {
			fmt.Fprintf(cmd.OutOrStdout(), "image %s already present\n", srv.Image)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "pulling %s\n", srv.Image)
			if err := cli.Pull(cmd.Context(), srv.Image, cmd.OutOrStdout()); err != nil {
				return err
			}
		}

		name := "anvil-" + id
		found, err := findNamed(cmd.Context(), cli, name)
		if err != nil {
			return err
		}

		dataDir := filepath.Join(paths.DataDir, "data", id)
		if err := os.MkdirAll(dataDir, 0o700); err != nil {
			return err
		}

		if found.ID == "" {
			mem, err := parseMemory(srv.Memory)
			if err != nil {
				return err
			}
			bind := dataDir + ":/data"
			if detected.SELinux == "enforcing" &&
				(detected.Mode == engine.ModeRootlessPodman ||
					detected.Mode == engine.ModeSystemPodman) {
				bind += ":Z"
			}
			publish := !srv.Public
			if srv.PublishToHost != nil {
				publish = *srv.PublishToHost
			}
			ports := []engine.Port{}
			for _, p := range srv.Ports {
				pt := engine.Port{Proto: "tcp", Container: p}
				if srv.Game != "minecraft" {
					pt.Proto = "udp"
				}
				if publish {
					pt.Host = p
				}
				ports = append(ports, pt)
			}
			iso := cli.Isolation()
			if iso.EnvUID != nil && srv.Env != nil {
				if srv.Env["UID"] == "" {
					srv.Env["UID"] = strconv.Itoa(*iso.EnvUID)
				}
				if srv.Env["GID"] == "" && iso.EnvGID != nil {
					srv.Env["GID"] = strconv.Itoa(*iso.EnvGID)
				}
			}
			if srv.Env == nil {
				srv.Env = map[string]string{}
			}
			if srv.Memory != "" && srv.Env["MEMORY"] == "" {
				srv.Env["MEMORY"] = srv.Memory
			}
			cid, err := cli.Create(cmd.Context(), engine.ContainerSpec{
				Name:        name,
				Image:       srv.Image,
				Env:         srv.Env,
				Binds:       []string{bind},
				Network:     "anvil",
				Ports:       ports,
				MemoryBytes: mem,
				Isolation:   iso,
				Labels: map[string]string{
					"io.anvil.managed": "true",
					"io.anvil.role":    "game",
					"io.anvil.server":  id,
					"io.anvil.game":    srv.Game,
					"io.anvil.version": "1",
				},
			})
			if err != nil {
				return err
			}
			found.ID = cid
		}

		if err := cli.Start(cmd.Context(), found.ID); err != nil {
			return err
		}
		ins, err := cli.Inspect(cmd.Context(), found.ID)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", id, srv.Game, ins.State)
		return nil
	},
}

func findNamed(ctx context.Context, cli *engine.Client, name string) (engine.Container, error) {
	list, err := cli.ListByLabel(ctx, map[string]string{"io.anvil.managed": "true"})
	if err != nil {
		return engine.Container{}, err
	}
	for _, c := range list {
		if c.Name == name {
			return c, nil
		}
	}
	return engine.Container{}, nil
}
