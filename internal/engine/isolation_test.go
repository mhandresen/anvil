package engine

import "testing"

func TestIsolationRootlessPodmanKeepsImageUID(t *testing.T) {
	id := TemplateIdentity{ImageUID: 1000, ImageGID: 1000}

	got := IsolationFor(ModeRootlessPodman, id, 1000, 1000)
	if got.UsernsMode != "keep-id:uid=1000,gid=1000" {
		t.Fatalf("userns=%q", got.UsernsMode)
	}
	if got.EnvUID != nil {
		t.Fatal("skal ikke sette PUID på rootless podman")
	}

	got = IsolationFor(ModeRootlessPodman, id, 1001, 1001)
	if got.UsernsMode != "keep-id:uid=1000,gid=1000" {
		t.Fatalf("host 1001 skal likevel bruke image-uid, fikk %q", got.UsernsMode)
	}
	if got.UsernsMode == "keep-id" {
		t.Fatal("bart keep-id er forbudt")
	}
}

func TestIsolationRootlessDocker(t *testing.T) {
	id := TemplateIdentity{ImageUID: 1000, ImageGID: 1000}
	got := IsolationFor(ModeRootlessDocker, id, 1000, 1000)
	if got.UsernsMode != "" {
		t.Fatalf("uventet userns %q", got.UsernsMode)
	}
	if got.ContainerUser != "0:0" {
		t.Fatalf("user=%q", got.ContainerUser)
	}
}

func TestIsolationSystemDocker(t *testing.T) {
	id := TemplateIdentity{ImageUID: 1000, ImageGID: 1000}
	got := IsolationFor(ModeSystemDocker, id, 1000, 1000)
	if got.ContainerUser != "1000:1000" {
		t.Fatalf("user=%q", got.ContainerUser)
	}
	if got.EnvUID == nil || *got.EnvUID != 1000 {
		t.Fatal("forventet PUID=host")
	}
}