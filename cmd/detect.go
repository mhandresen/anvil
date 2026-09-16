package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mhandresen/anvil/internal/engine"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect Podman or Docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := engine.Detect(cmd.Context(), os.Getuid(), envMap())
		if err == nil && d.Socket != "" {
			iso := engine.IsolationFromMode(d.Mode, os.Getuid(), os.Getgid())
			cli, api, cerr := engine.Connect(cmd.Context(), d.Socket, d.Kind, iso)
			if cerr != nil {
				err = cerr
			} else {
				d.API = api
				_ = cli.Close()
			}
		}

		payload := detectOutput{
			Engine:  none(string(d.Kind)),
			Mode:    none(string(d.Mode)),
			Socket:  none(d.Socket),
			API:     none(d.API),
			OK:      err == nil && d.Kind != "",
			Linger:  d.Linger,
			SELinux: d.SELinux,
		}

		if flagJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			if encErr := enc.Encode(payload); encErr != nil {
				return encErr
			}
		} else {
			printDetect(payload, d.Tried)
		}

		return err
	},
}

type detectOutput struct {
	Engine  string `json:"engine"`
	Mode    string `json:"mode"`
	Socket  string `json:"socket"`
	API     string `json:"api"`
	OK      bool   `json:"ok"`
	Linger  *bool  `json:"linger"`
	SELinux string `json:"selinux"`
}

func printDetect(p detectOutput, tried []string) {
	fmt.Printf("engine:    %s\n", p.Engine)
	fmt.Printf("mode:      %s\n", p.Mode)
	fmt.Printf("socket:    %s\n", p.Socket)
	fmt.Printf("api:       %s\n", p.API)
	fmt.Printf("ok:        %v\n", p.OK)
	if p.SELinux != "" {
		fmt.Printf("selinux:   %s\n", p.SELinux)
	}
	if p.Linger != nil {
		fmt.Printf("linger:    %v\n", *p.Linger)
	}
	if !p.OK {
		if p.Socket != "(none)" && p.Socket != "" {
			fmt.Printf("\nsocket found but engine API is unreachable\n")
		} else {
			fmt.Printf("\nno container engine found\n")
			if len(tried) > 0 {
				fmt.Printf("tried:\n")
				for _, c := range tried {
					fmt.Printf("  %s\n", c)
				}
			}
			fmt.Printf("on Linux: systemctl --user enable --now podman.socket\n")
		}
	}
}

func none(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

func envMap() map[string]string {
	out := make(map[string]string, 16)
	for _, e := range os.Environ() {
		k, v, ok := strings.Cut(e, "=")
		if ok {
			out[k] = v
		}
	}
	return out
}