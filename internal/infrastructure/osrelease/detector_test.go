package osrelease_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/osrelease"
)

func makeReader(content string) func() (io.Reader, error) {
	return func() (io.Reader, error) {
		return strings.NewReader(content), nil
	}
}

func defaultDE() domain.DesktopEnvironment { return domain.Other }

func TestDetector_Debian12(t *testing.T) {
	content := "PRETTY_NAME=\"Debian GNU/Linux 12 (bookworm)\"\nID=debian\nVERSION_ID=\"12\"\n"
	det := osrelease.NewDetector(makeReader(content), noEnv, defaultDE)
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
	det := osrelease.NewDetector(makeReader(content), noEnv, defaultDE)
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
	det := osrelease.NewDetector(makeReader(content), noEnv, defaultDE)
	_, err := det.Detect()
	if err == nil {
		t.Fatal("expected error for missing VERSION_ID")
	}
}

func TestDetector_MissingID(t *testing.T) {
	content := "VERSION_ID=\"12\"\n"
	det := osrelease.NewDetector(makeReader(content), noEnv, defaultDE)
	_, err := det.Detect()
	if err == nil {
		t.Fatal("expected error for missing ID")
	}
}

func TestDetector_ReaderError(t *testing.T) {
	det := osrelease.NewDetector(func() (io.Reader, error) {
		return nil, errors.New("file not found: /etc/os-release")
	}, noEnv, defaultDE)
	_, err := det.Detect()
	if err == nil {
		t.Fatal("expected error when reader fails")
	}
}

func noEnv(string) string { return "" }

func TestDetector_DesktopEnvironment(t *testing.T) {
	content := "ID=debian\nVERSION_ID=\"12\"\n"
	
	t.Run("Detect KDE via env", func(t *testing.T) {
		mockEnv := func(k string) string {
			if k == "XDG_CURRENT_DESKTOP" {
				return "KDE"
			}
			return ""
		}
		det := osrelease.NewDetector(makeReader(content), mockEnv, defaultDE)
		info, err := det.Detect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.DesktopEnvironment != domain.KDE {
			t.Errorf("expected DE=kde, got %q", info.DesktopEnvironment)
		}
	})

	t.Run("Detect GNOME via env", func(t *testing.T) {
		mockEnv := func(k string) string {
			if k == "XDG_CURRENT_DESKTOP" {
				return "GNOME"
			}
			return ""
		}
		det := osrelease.NewDetector(makeReader(content), mockEnv, defaultDE)
		info, err := det.Detect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.DesktopEnvironment != domain.GNOME {
			t.Errorf("expected DE=gnome, got %q", info.DesktopEnvironment)
		}
	})

	t.Run("Fallback to deDetectorFn", func(t *testing.T) {
		det := osrelease.NewDetector(makeReader(content), noEnv, func() domain.DesktopEnvironment {
			return domain.GNOME
		})
		info, err := det.Detect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.DesktopEnvironment != domain.GNOME {
			t.Errorf("expected DE=gnome from fallback, got %q", info.DesktopEnvironment)
		}
	})
}
