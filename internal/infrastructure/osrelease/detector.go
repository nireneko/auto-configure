package osrelease

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/so-install/internal/core/domain"
)

// Detector reads and parses /etc/os-release.
type Detector struct {
	readerFn func() (io.Reader, error)
}

// NewDetector creates a Detector with a custom reader function (for testability).
func NewDetector(readerFn func() (io.Reader, error)) *Detector {
	return &Detector{readerFn: readerFn}
}

// NewDefaultDetector creates a Detector that reads the real /etc/os-release.
func NewDefaultDetector() *Detector {
	return NewDetector(func() (io.Reader, error) {
		return os.Open("/etc/os-release")
	})
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

	return &domain.OSInfo{ID: id, VersionID: versionID}, nil
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
