package engine

import (
	"context"
	"fmt"
	"io"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type Client struct {
	kind Kind
	iso  Isolation
	cli  *client.Client
	host string
}

func Connect(ctx context.Context, socket string, kind Kind, iso Isolation) (*Client, string, error) {
	cli, err := client.New(client.WithHost(socket))
	if err != nil {
		return nil, "", err
	}
	ping, err := cli.Ping(ctx, client.PingOptions{})
	if err != nil {
		_ = cli.Close()
		return nil, "", err
	}
	return &Client{kind: kind, iso: iso, cli: cli, host: socket}, ping.APIVersion, nil
}

func (c *Client) Close() error {
	if c == nil || c.cli == nil {
		return nil
	}
	return c.cli.Close()
}

func (c *Client) Kind() Kind           { return c.kind }
func (c *Client) Isolation() Isolation { return c.iso }

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx, client.PingOptions{})
	return err
}

func (c *Client) EnsureNetwork(ctx context.Context, spec NetworkSpec) (Network, error) {
	name := spec.Name
	if name == "" {
		name = "anvil"
	}
	list, err := c.cli.NetworkList(ctx, client.NetworkListOptions{
		Filters: client.Filters{"name": {"true": false}},
	})
	if err != nil {
		return Network{}, err
	}
	// Filters values are map[string]bool; match by iterating.
	_ = list
	list, err = c.cli.NetworkList(ctx, client.NetworkListOptions{
		Filters: client.Filters{"name": {name: true}},
	})
	if err != nil {
		return Network{}, err
	}
	for _, n := range list.Items {
		if n.Name == name {
			return Network{Name: name, Subnet: spec.Subnet}, nil
		}
	}

	opts := client.NetworkCreateOptions{
		Driver: "bridge",
		Labels: map[string]string{
			"io.anvil.managed": "true",
			"io.anvil.role":    "network",
			"io.anvil.scope":   "host",
		},
	}
	if spec.Subnet != "" {
		pfx, err := netip.ParsePrefix(spec.Subnet)
		if err != nil {
			return Network{}, fmt.Errorf("invalid subnet %q: %w", spec.Subnet, err)
		}
		opts.IPAM = &network.IPAM{
			Config: []network.IPAMConfig{{Subnet: pfx}},
		}
	}
	if _, err := c.cli.NetworkCreate(ctx, name, opts); err != nil {
		return Network{}, err
	}
	return Network{Name: name, Subnet: spec.Subnet}, nil
}

func (c *Client) Pull(ctx context.Context, ref string, w io.Writer) error {
	resp, err := c.cli.ImagePull(ctx, ref, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()
	if w == nil {
		w = io.Discard
	}
	_, err = io.Copy(w, resp)
	return err
}

func (c *Client) Create(ctx context.Context, spec ContainerSpec) (ContainerID, error) {
	exposed := network.PortSet{}
	bindings := network.PortMap{}
	for _, p := range spec.Ports {
		proto := p.Proto
		if proto == "" {
			proto = "tcp"
		}
		port, err := network.ParsePort(fmt.Sprintf("%d/%s", p.Container, proto))
		if err != nil {
			return "", err
		}
		exposed[port] = struct{}{}
		if p.Host > 0 {
			bindings[port] = []network.PortBinding{{
				HostIP:   netip.IPv4Unspecified(),
				HostPort: strconv.Itoa(p.Host),
			}}
		}
	}

	env := make([]string, 0, len(spec.Env))
	for k, v := range spec.Env {
		env = append(env, k+"="+v)
	}

	hostCfg := &container.HostConfig{
		Binds:         spec.Binds,
		PortBindings:  bindings,
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
	}
	if spec.MemoryBytes > 0 {
		hostCfg.Memory = spec.MemoryBytes
	}
	if spec.NanoCPUs > 0 {
		hostCfg.NanoCPUs = spec.NanoCPUs
	}
	if spec.Isolation.UsernsMode != "" {
		hostCfg.UsernsMode = container.UsernsMode(spec.Isolation.UsernsMode)
	}

	cfg := &container.Config{
		Image:        spec.Image,
		Env:          env,
		Labels:       spec.Labels,
		ExposedPorts: exposed,
	}
	if spec.Isolation.ContainerUser != "" {
		cfg.User = spec.Isolation.ContainerUser
	}

	var ep *network.NetworkingConfig
	if spec.Network != "" {
		epc := &network.EndpointSettings{}
		if spec.IP != "" {
			addr, err := netip.ParseAddr(spec.IP)
			if err != nil {
				return "", fmt.Errorf("invalid ip %q: %w", spec.IP, err)
			}
			epc.IPAMConfig = &network.EndpointIPAMConfig{IPv4Address: addr}
		}
		ep = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				spec.Network: epc,
			},
		}
	}

	resp, err := c.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:             spec.Name,
		Config:           cfg,
		HostConfig:       hostCfg,
		NetworkingConfig: ep,
	})
	if err != nil {
		return "", err
	}
	return ContainerID(resp.ID), nil
}

func (c *Client) Start(ctx context.Context, id ContainerID) error {
	_, err := c.cli.ContainerStart(ctx, string(id), client.ContainerStartOptions{})
	return err
}

func (c *Client) Stop(ctx context.Context, id ContainerID, timeout time.Duration) error {
	sec := int(timeout.Seconds())
	_, err := c.cli.ContainerStop(ctx, string(id), client.ContainerStopOptions{Timeout: &sec})
	return err
}

func (c *Client) Remove(ctx context.Context, id ContainerID, volumes bool) error {
	_, err := c.cli.ContainerRemove(ctx, string(id), client.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: volumes,
	})
	return err
}

func (c *Client) Inspect(ctx context.Context, id ContainerID) (Container, error) {
	ins, err := c.cli.ContainerInspect(ctx, string(id), client.ContainerInspectOptions{})
	if err != nil {
		return Container{}, err
	}
	ctr := ins.Container
	out := Container{
		ID:   ContainerID(ctr.ID),
		Name: strings.TrimPrefix(ctr.Name, "/"),
	}
	if ctr.State != nil {
		out.State = string(ctr.State.Status)
	}
	if ctr.NetworkSettings != nil {
		for _, n := range ctr.NetworkSettings.Networks {
			if n != nil && n.IPAddress.IsValid() {
				out.IP = n.IPAddress.String()
				break
			}
		}
	}
	return out, nil
}

func (c *Client) ListByLabel(ctx context.Context, labels map[string]string) ([]Container, error) {
	set := map[string]bool{}
	for k, v := range labels {
		set[k+"="+v] = true
	}
	list, err := c.cli.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: client.Filters{"label": set},
	})
	if err != nil {
		return nil, err
	}
	out := make([]Container, 0, len(list.Items))
	for _, it := range list.Items {
		name := ""
		if len(it.Names) > 0 {
			name = strings.TrimPrefix(it.Names[0], "/")
		}
		out = append(out, Container{
			ID:    ContainerID(it.ID),
			Name:  name,
			State: string(it.State),
		})
	}
	return out, nil
}

func (c *Client) Logs(ctx context.Context, id ContainerID, follow bool) (io.ReadCloser, error) {
	return c.cli.ContainerLogs(ctx, string(id), client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
	})
}

func (c *Client) Stats(ctx context.Context, id ContainerID) (Stats, error) {
	_ = ctx
	_ = id
	return Stats{}, ErrNotImplemented
}

func (c *Client) ImageExists(ctx context.Context, ref string) bool {
	_, err := c.cli.ImageInspect(ctx, ref)
	return err == nil
}