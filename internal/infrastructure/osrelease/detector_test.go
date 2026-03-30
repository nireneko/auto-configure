package osrelease_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/so-install/internal/infrastructure/osrelease"
)

func makeReader(content string) func() (io.Reader, error) {
	return func() (io.Reader, error) {
		return strings.NewReader(content), nil
	}
}

func TestDetector_Debian12(t *testing.T) {
	content := "PRETTY_NAME=\"Debian GNU/Linux 12 (bookworm)\"\nID=debian\nVERSION_ID=\"12\"\n"
	det := osrelease.NewDetector(makeReader(content))
	info, err := det.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "debian" {
		t.Errorf("expected ID=debian, got %q", info.ID)
	}
	if info.VersionID != "12" {
		t.Errorf("expected VersionID=12, got %q", info.VersionID)
	}
}

func TestDetector_Debian13(t *testing.T) {
	content := "ID=debian\nVERSION_ID=\"13\"\n"
	det := osrelease.NewDetector(makeReader(content))
	info, err := det.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "debian" || info.VersionID != "13" {
		t.Errorf("wrong info: %+v", info)
	}
}

func TestDetector_MissingVersionID(t *testing.T) {
	content := "ID=debian\n"
	det := osrelease.NewDetector(makeReader(content))
	_, err := det.Detect()
	if err == nil {
		t.Fatal("expected error for missing VERSION_ID")
	}
}

func TestDetector_MissingID(t *testing.T) {
	content := "VERSION_ID=\"12\"\n"
	det := osrelease.NewDetector(makeReader(content))
	_, err := det.Detect()
	if err == nil {
		t.Fatal("expected error for missing ID")
	}
}

func TestDetector_ReaderError(t *testing.T) {
	det := osrelease.NewDetector(func() (io.Reader, error) {
		return nil, errors.New("file not found: /etc/os-release")
	})
	_, err := det.Detect()
	if err == nil {
		t.Fatal("expected error when reader fails")
	}
}
