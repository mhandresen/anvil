package engine

type Relabel string

const (
	RelabelNone    Relabel = "none"
	RelabelPrivate Relabel = "Z"
	RelabelShared  Relabel = "z"
)

type Isolation struct {
	Mode          Mode
	UsernsMode    string
	ContainerUser string
	EnvUID        *int
	EnvGID        *int
	MountRelabel  Relabel
}

type TemplateIdentity struct {
	ImageUID int
	ImageGID int
}

type NetworkSpec struct {
	Name   string
	Subnet string
}

type Network struct {
	Name   string
	Subnet string
}

type ContainerID string

type Port struct {
	Proto     string
	Container int
	Host      int
}

type ContainerSpec struct {
	Name        string
	Image       string
	Labels      map[string]string
	Env         map[string]string
	Binds       []string
	Network     string
	IP          string
	Ports       []Port
	MemoryBytes int64
	NanoCPUs    int64
	Isolation   Isolation
}

type Container struct {
	ID    ContainerID
	Name  string
	State string
	IP    string
}

type Stats struct {
	CPUPercent float64
	MemoryRSS  uint64
}