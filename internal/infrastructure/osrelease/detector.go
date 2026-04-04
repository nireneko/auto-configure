package osrelease

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/so-install/internal/core/domain"
)

// Detector reads and parses /etc/os-release and detects the Desktop Environment.
type Detector struct {
	readerFn     func() (io.Reader, error)
	envFn        func(string) string
	deDetectorFn func() domain.DesktopEnvironment
}

// NewDetector creates a Detector with custom dependencies.
func NewDetector(
	readerFn func() (io.Reader, error),
	envFn func(string) string,
	deDetectorFn func() domain.DesktopEnvironment,
) *Detector {
	return &Detector{
		readerFn:     readerFn,
		envFn:        envFn,
		deDetectorFn: deDetectorFn,
	}
}

// NewDefaultDetector creates a Detector with real system dependencies.
func NewDefaultDetector() *Detector {
	return NewDetector(
		func() (io.Reader, error) {
			return os.Open("/etc/os-release")
		},
		os.Getenv,
		detectDesktopEnvironment,
	)
}

var _ domain.OSDetector = (*Detector)(nil)

// Detect parses /etc/os-release and returns OSInfo.
func (d *Detector) Detect() (*domain.OSInfo, error) {
	r, err := d.readerFn()
	if err != nil {
		return nil, fmt.Errorf("cannot read OS information: %w", err)
	}

	kvs := parseKeyValues(r)

	id, ok := kvs["ID"]
	if !ok || id == "" {
		return nil, fmt.Errorf("ID field not found in os-release")
	}

	versionID, ok := kvs["VERSION_ID"]
	if !ok || versionID == "" {
		return nil, fmt.Errorf("VERSION_ID field not found in os-release")
	}

	de := d.detectDE()
	isWayland := d.envFn("XDG_SESSION_TYPE") == "wayland" || d.envFn("WAYLAND_DISPLAY") != ""

	return &domain.OSInfo{
		ID:                 id,
		VersionID:          versionID,
		DesktopEnvironment: de,
		IsWayland:          isWayland,
	}, nil
}

func (d *Detector) detectDE() domain.DesktopEnvironment {
	xdg := strings.ToLower(d.envFn("XDG_CURRENT_DESKTOP"))
	if strings.Contains(xdg, "kde") {
		return domain.KDE
	}
	if strings.Contains(xdg, "gnome") {
		return domain.GNOME
	}

	// Fallback to more expensive detection
	return d.deDetectorFn()
}

func detectDesktopEnvironment() domain.DesktopEnvironment {
	// 1. Check for running processes
	if isProcessRunning("gnome-shell") {
		return domain.GNOME
	}
	if isProcessRunning("plasmashell") || isProcessRunning("kwin_x11") || isProcessRunning("kwin_wayland") {
		return domain.KDE
	}

	// 2. Check for installed packages (last resort)
	if isPackageInstalled("gnome-shell") {
		return domain.GNOME
	}
	if isPackageInstalled("plasma-desktop") {
		return domain.KDE
	}

	return domain.Other
}

func isProcessRunning(name string) bool {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		commPath := fmt.Sprintf("/proc/%s/comm", f.Name())
		content, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(content)) == name {
			return true
		}
	}
	return false
}

func isPackageInstalled(name string) bool {
	// Simplistic check for Debian-based systems
	_, err := os.Stat(fmt.Sprintf("/var/lib/dpkg/info/%s.list", name))
	return err == nil
}

// parseKeyValues parses key=value lines, stripping surrounding quotes.
func parseKeyValues(r io.Reader) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Strip surrounding quotes
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		result[key] = val
	}
	return result
}
