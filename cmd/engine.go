package cmd

import (
	"context"
	"os"

	"github.com/mhandresen/anvil/internal/engine"
)

func openEngine(ctx context.Context) (*engine.Client, engine.Detected, error) {
	d, err := engine.Detect(ctx, os.Getuid(), envMap())
	if err != nil {
		return nil, d, err
	}
	iso := engine.IsolationFromMode(d.Mode, os.Getuid(), os.Getgid())
	cli, api, err := engine.Connect(ctx, d.Socket, d.Kind, iso)
	if err != nil {
		return nil, d, err
	}
	d.API = api
	return cli, d, nil
}