package engine

import (
	"context"
	"io"
	"time"
)

// FakeEngine er testdobbel. Ingen daemon.
type FakeEngine struct {
	K   Kind
	Iso Isolation
}

func (f *FakeEngine) Kind() Kind             { return f.K }
func (f *FakeEngine) Isolation() Isolation   { return f.Iso }
func (f *FakeEngine) Ping(context.Context) error { return nil }

func (f *FakeEngine) EnsureNetwork(context.Context, NetworkSpec) (Network, error) {
	return Network{}, ErrNotImplemented
}
func (f *FakeEngine) Pull(context.Context, string, io.Writer) error { return ErrNotImplemented }
func (f *FakeEngine) Create(context.Context, ContainerSpec) (ContainerID, error) {
	return "", ErrNotImplemented
}
func (f *FakeEngine) Start(context.Context, ContainerID) error { return ErrNotImplemented }
func (f *FakeEngine) Stop(context.Context, ContainerID, time.Duration) error {
	return ErrNotImplemented
}
func (f *FakeEngine) Remove(context.Context, ContainerID, bool) error { return ErrNotImplemented }
func (f *FakeEngine) Inspect(context.Context, ContainerID) (Container, error) {
	return Container{}, ErrNotImplemented
}
func (f *FakeEngine) ListByLabel(context.Context, map[string]string) ([]Container, error) {
	return nil, ErrNotImplemented
}
func (f *FakeEngine) Logs(context.Context, ContainerID, bool) (io.ReadCloser, error) {
	return nil, ErrNotImplemented
}
func (f *FakeEngine) Stats(context.Context, ContainerID) (Stats, error) {
	return Stats{}, ErrNotImplemented
}