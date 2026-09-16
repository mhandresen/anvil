package engine

import "fmt"

func IsolationFor(mode Mode, id TemplateIdentity, hostUID, hostGID int) Isolation {
	iso := Isolation{
		Mode:         mode,
		MountRelabel: RelabelNone,
	}

	switch mode {
	case ModeRootlessPodman:
		iso.UsernsMode = fmt.Sprintf("keep-id:uid=%d,gid=%d", id.ImageUID, id.ImageGID)
	case ModeRootlessDocker:
		iso.ContainerUser = "0:0"
	case ModeSystemDocker, ModeSystemPodman:
		iso.ContainerUser = fmt.Sprintf("%d:%d", hostUID, hostGID)
		u, g := hostUID, hostGID
		iso.EnvUID = &u
		iso.EnvGID = &g
	}

	return iso
}