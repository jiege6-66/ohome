package service

import (
	"errors"
	"testing"
)

func TestResolveQuarkDirectLinkForEntry_VideoUsesTVLink(t *testing.T) {
	entry := quarkDriveFile{Category: 1, Size: 1024, File: true}
	videoCalled := false
	downloadCalled := false

	link, err := resolveQuarkDirectLinkForEntry(
		entry,
		func() (quarkFileLink, error) {
			videoCalled = true
			return quarkFileLink{URL: "tv-link"}, nil
		},
		func() (quarkFileLink, error) {
			downloadCalled = true
			return quarkFileLink{URL: "download-link"}, nil
		},
	)
	if err != nil {
		t.Fatalf("resolveQuarkDirectLinkForEntry() error = %v", err)
	}
	if !videoCalled {
		t.Fatalf("expected video resolver to be called")
	}
	if downloadCalled {
		t.Fatalf("did not expect download resolver to be called")
	}
	if link.URL != "tv-link" {
		t.Fatalf("expected tv link, got %q", link.URL)
	}
}

func TestResolveQuarkDirectLinkForEntry_NonVideoUsesDownloadLink(t *testing.T) {
	entry := quarkDriveFile{Category: 2, Size: 1024, File: true}
	videoCalled := false
	downloadCalled := false

	link, err := resolveQuarkDirectLinkForEntry(
		entry,
		func() (quarkFileLink, error) {
			videoCalled = true
			return quarkFileLink{URL: "tv-link"}, nil
		},
		func() (quarkFileLink, error) {
			downloadCalled = true
			return quarkFileLink{URL: "download-link"}, nil
		},
	)
	if err != nil {
		t.Fatalf("resolveQuarkDirectLinkForEntry() error = %v", err)
	}
	if videoCalled {
		t.Fatalf("did not expect video resolver to be called")
	}
	if !downloadCalled {
		t.Fatalf("expected download resolver to be called")
	}
	if link.URL != "download-link" {
		t.Fatalf("expected download link, got %q", link.URL)
	}
}

func TestResolveQuarkDirectLinkForEntry_VideoErrorReturned(t *testing.T) {
	entry := quarkDriveFile{Category: 1, Size: 1024, File: true}
	wantErr := errors.New("tv login missing")

	_, err := resolveQuarkDirectLinkForEntry(
		entry,
		func() (quarkFileLink, error) {
			return quarkFileLink{}, wantErr
		},
		func() (quarkFileLink, error) {
			return quarkFileLink{URL: "download-link"}, nil
		},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}
