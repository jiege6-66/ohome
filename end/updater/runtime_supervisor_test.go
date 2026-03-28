package updater

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"ohome/buildinfo"
)

func TestBootstrapCurrentReleaseSeedsEmbeddedBinaryOnFirstBoot(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("launcher bootstrap uses symlinks and is only exercised in Linux containers")
	}

	originalRuntimeVersion := buildinfo.RuntimeVersion
	t.Cleanup(func() {
		buildinfo.RuntimeVersion = originalRuntimeVersion
	})
	buildinfo.RuntimeVersion = "runtime-v2026.03.1"

	rootDir := t.TempDir()
	embeddedServerPath := filepath.Join(rootDir, "embedded", "server")
	if err := os.MkdirAll(filepath.Dir(embeddedServerPath), 0o755); err != nil {
		t.Fatalf("mkdir embedded dir: %v", err)
	}
	if err := os.WriteFile(embeddedServerPath, []byte("embedded-server"), 0o755); err != nil {
		t.Fatalf("write embedded server: %v", err)
	}

	store := &Store{rootDir: filepath.Join(rootDir, "data", "update")}
	supervisor := NewRuntimeSupervisor(store, rootDir, embeddedServerPath)

	version, releasePath, err := supervisor.BootstrapCurrentRelease("v1.2.3")
	if err != nil {
		t.Fatalf("BootstrapCurrentRelease returned error: %v", err)
	}
	if version != "v1.2.3" {
		t.Fatalf("unexpected version: %s", version)
	}

	linkTarget, err := resolveCurrentReleasePath(store.CurrentReleaseLink())
	if err != nil {
		t.Fatalf("resolve current release link: %v", err)
	}
	if linkTarget != filepath.Clean(releasePath) {
		t.Fatalf("current release link=%s, want %s", linkTarget, filepath.Clean(releasePath))
	}

	payload, err := os.ReadFile(filepath.Join(releasePath, "server"))
	if err != nil {
		t.Fatalf("read seeded server: %v", err)
	}
	if string(payload) != "embedded-server" {
		t.Fatalf("unexpected seeded server payload: %s", string(payload))
	}
}

func TestBootstrapCurrentReleaseSeedsEmbeddedBinaryOnRuntimeChange(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("launcher bootstrap uses symlinks and is only exercised in Linux containers")
	}

	originalRuntimeVersion := buildinfo.RuntimeVersion
	t.Cleanup(func() {
		buildinfo.RuntimeVersion = originalRuntimeVersion
	})
	buildinfo.RuntimeVersion = "runtime-v2026.03.2"

	rootDir := t.TempDir()
	embeddedServerPath := filepath.Join(rootDir, "embedded", "server")
	if err := os.MkdirAll(filepath.Dir(embeddedServerPath), 0o755); err != nil {
		t.Fatalf("mkdir embedded dir: %v", err)
	}
	if err := os.WriteFile(embeddedServerPath, []byte("runtime-embedded"), 0o755); err != nil {
		t.Fatalf("write embedded server: %v", err)
	}

	store := &Store{rootDir: filepath.Join(rootDir, "data", "update")}
	currentReleasePath := filepath.Join(store.ReleasesDir(), "v9.9.9")
	if err := os.MkdirAll(currentReleasePath, 0o755); err != nil {
		t.Fatalf("mkdir current release: %v", err)
	}
	if err := os.WriteFile(filepath.Join(currentReleasePath, "server"), []byte("current-server"), 0o755); err != nil {
		t.Fatalf("write current release server: %v", err)
	}
	if err := setCurrentReleaseLink(store.CurrentReleaseLink(), currentReleasePath); err != nil {
		t.Fatalf("set current release link: %v", err)
	}
	if err := store.SaveState(&State{
		CurrentVersion:     "v9.9.9",
		CurrentReleasePath: currentReleasePath,
		RuntimeVersion:     "runtime-v2026.03.1",
	}); err != nil {
		t.Fatalf("save state: %v", err)
	}

	supervisor := NewRuntimeSupervisor(store, rootDir, embeddedServerPath)
	version, releasePath, err := supervisor.BootstrapCurrentRelease("v1.2.3")
	if err != nil {
		t.Fatalf("BootstrapCurrentRelease returned error: %v", err)
	}
	if version != "v9.9.9" {
		t.Fatalf("unexpected active version: %s", version)
	}
	if filepath.Clean(releasePath) != filepath.Clean(currentReleasePath) {
		t.Fatalf("active release path=%s, want %s", releasePath, currentReleasePath)
	}

	seededPath := filepath.Join(store.ReleasesDir(), "v1.2.3", "server")
	payload, err := os.ReadFile(seededPath)
	if err != nil {
		t.Fatalf("read seeded runtime server: %v", err)
	}
	if string(payload) != "runtime-embedded" {
		t.Fatalf("unexpected seeded runtime payload: %s", string(payload))
	}

	linkTarget, err := resolveCurrentReleasePath(store.CurrentReleaseLink())
	if err != nil {
		t.Fatalf("resolve current release link: %v", err)
	}
	if filepath.Clean(linkTarget) != filepath.Clean(currentReleasePath) {
		t.Fatalf("current release link changed to %s, want %s", linkTarget, currentReleasePath)
	}
}
