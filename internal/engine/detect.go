package engine

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Detected struct {
	Kind    Kind
	Mode    Mode
	Socket  string
	API     string
	Linger  *bool
	SELinux string
	Tried   []string
}

func Candidates(uid int, env map[string]string) []string {
	var out []string
	seen := map[string]bool{}
	add := func(s string) {
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		out = append(out, s)
	}

	add(env["ANVIL_HOST"])
	add(env["DOCKER_HOST"])
	add(env["CONTAINER_HOST"])
	add(fmt.Sprintf("unix:///run/user/%d/podman/podman.sock", uid))
	add(fmt.Sprintf("unix:///run/user/%d/docker.sock", uid))
	add("unix:///var/run/docker.sock")
	add("unix:///run/podman/podman.sock")
	return out
}

func Detect(ctx context.Context, uid int, env map[string]string) (Detected, error) {
	return DetectFrom(ctx, Candidates(uid, env))
}

func DetectFrom(ctx context.Context, candidates []string) (Detected, error) {
	d := Detected{
		Tried:   candidates,
		Linger:  probeLinger(),
		SELinux: probeSELinux(),
	}

	for _, cand := range candidates {
		if err := ctxErr(ctx); err != nil {
			return d, err
		}
		if !socketExists(cand) {
			continue
		}
		d.Socket = cand
		d.Kind, d.Mode = guessKindMode(cand)
		return d, nil
	}
	return d, ErrNoEngine
}

func ctxErr(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}

func socketExists(host string) bool {
	path := strings.TrimPrefix(host, "unix://")
	if strings.Contains(path, "://") {
		return false
	}
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.Mode()&os.ModeSocket != 0 || st.Mode().IsRegular()
}

func guessKindMode(socket string) (Kind, Mode) {
	podman := strings.Contains(socket, "podman")
	userRun := strings.Contains(socket, "/run/user/")
	switch {
	case podman && userRun:
		return KindPodman, ModeRootlessPodman
	case podman:
		return KindPodman, ModeSystemPodman
	case userRun:
		return KindDocker, ModeRootlessDocker
	default:
		return KindDocker, ModeSystemDocker
	}
}

func probeLinger() *bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	user := lastPath(home)
	if user == "" {
		return nil
	}
	p := "/var/lib/systemd/linger/" + user
	if _, err := os.Stat(p); err == nil {
		v := true
		return &v
	}
	if _, err := os.Stat("/var/lib/systemd/linger"); err == nil {
		v := false
		return &v
	}
	return nil
}

func lastPath(home string) string {
	home = strings.TrimRight(home, "/")
	i := strings.LastIndex(home, "/")
	if i < 0 || i == len(home)-1 {
		return ""
	}
	return home[i+1:]
}

func probeSELinux() string {
	b, err := os.ReadFile("/sys/fs/selinux/enforce")
	if err != nil {
		return ""
	}
	switch strings.TrimSpace(string(b)) {
	case "1":
		return "enforcing"
	case "0":
		return "permissive"
	default:
		return ""
	}
}

func IsolationFromMode(mode Mode, hostUID, hostGID int) Isolation {
	return IsolationFor(mode, TemplateIdentity{ImageUID: 1000, ImageGID: 1000}, hostUID, hostGID)
}