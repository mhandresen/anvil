package engine

type Kind string

const (
	KindPodman Kind = "podman"
	KindDocker Kind = "docker"
)

func (k Kind) String() string { return string(k) }

type Mode string

const (
	ModeRootlessPodman Mode = "rootless-podman"
	ModeRootlessDocker Mode = "rootless-docker"
	ModeSystemDocker   Mode = "system-docker"
	ModeSystemPodman   Mode = "system-podman"
)

func (m Mode) String() string { return string(m) }

func (m Mode) Rootless() bool {
	return m == ModeRootlessPodman || m == ModeRootlessDocker
}