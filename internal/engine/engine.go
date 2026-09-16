package engine

import (
	"context"
	"io"
	"time"
)

type Engine interface {
	Kind() Kind
	Isolation() Isolation
	Ping(ctx context.Context) error
	EnsureNetwork(ctx context.Context, spec NetworkSpec) (Network, error)
	Pull(ctx context.Context, ref string, w io.Writer) error
	Create(ctx context.Context, spec ContainerSpec) (ContainerID, error)
	Start(ctx context.Context, id ContainerID) error
	Stop(ctx context.Context, id ContainerID, timeout time.Duration) error
	Remove(ctx context.Context, id ContainerID, volumes bool) error
	Inspect(ctx context.Context, id ContainerID) (Container, error)
	ListByLabel(ctx context.Context, labels map[string]string) ([]Container, error)
	Logs(ctx context.Context, id ContainerID, follow bool) (io.ReadCloser, error)
	Stats(ctx context.Context, id ContainerID) (Stats, error)
}