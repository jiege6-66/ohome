package updater

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ohome/buildinfo"
)

func TestFetchManifest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"channel":"stable",
			"version":"v1.2.3",
			"releaseNotes":"Release v1.2.3",
			"publishedAt":"2026-03-28T00:00:00Z",
			"minRuntimeVersion":"runtime-v2026.03.1",
			"recommendedRuntimeVersion":"runtime-v2026.03.1",
			"artifacts":{
				"linux-amd64":{"url":"https://example.com/amd64.tar.gz","sha256":"abc"},
				"linux-arm64":{"url":"https://example.com/arm64.tar.gz","sha256":"def"}
			}
		}`))
	}))
	defer server.Close()

	manifest, err := FetchManifest(server.URL)
	if err != nil {
		t.Fatalf("FetchManifest returned error: %v", err)
	}
	if manifest.Version != "v1.2.3" {
		t.Fatalf("unexpected version: %s", manifest.Version)
	}
	if manifest.Artifacts["linux-amd64"].URL == "" {
		t.Fatalf("expected linux-amd64 artifact to be present")
	}
}

func TestValidateRuntimeCompatibility(t *testing.T) {
	t.Parallel()

	originalRuntimeVersion := buildinfo.RuntimeVersion
	t.Cleanup(func() {
		buildinfo.RuntimeVersion = originalRuntimeVersion
	})

	buildinfo.RuntimeVersion = "runtime-v2026.03.1"
	if err := validateRuntimeCompatibility(ServerManifest{MinRuntimeVersion: "runtime-v2026.03.1"}); err != nil {
		t.Fatalf("expected compatibility check to pass, got %v", err)
	}

	buildinfo.RuntimeVersion = "runtime-v2026.03.1"
	if err := validateRuntimeCompatibility(ServerManifest{MinRuntimeVersion: "runtime-v2026.03.2"}); err == nil {
		t.Fatalf("expected compatibility check to fail")
	}
}
