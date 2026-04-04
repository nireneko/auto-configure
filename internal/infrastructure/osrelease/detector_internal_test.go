package osrelease

import (
	"io"
	"strings"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestParseKeyValues(t *testing.T) {
	input := `
ID=debian
VERSION_ID="12"
# Comment
EMPTY=
NO_EQUALS
`
	kvs := parseKeyValues(strings.NewReader(input))
	assert.Equal(t, "debian", kvs["ID"])
	assert.Equal(t, "12", kvs["VERSION_ID"])
}

func TestNewDefaultDetector(t *testing.T) {
	d := NewDefaultDetector()
	assert.NotNil(t, d)
}

func TestIsProcessRunning(t *testing.T) {
	assert.False(t, IsProcessRunning("definitely-not-running-12345"))
}

func TestIsPackageInstalled(t *testing.T) {
	assert.False(t, isPackageInstalled("definitely-not-installed-12345"))
}

func TestDetectDesktopEnvironment_EdgeCases(t *testing.T) {
	de := DetectDesktopEnvironment()
	assert.Contains(t, []domain.DesktopEnvironment{domain.GNOME, domain.KDE, domain.Other}, de)
}

func TestDetector_Detect_Errors(t *testing.T) {
	det := NewDetector(func() (io.Reader, error) { return strings.NewReader("VERSION_ID=12"), nil }, func(string) string { return "" }, func() domain.DesktopEnvironment { return domain.Other })
	_, err := det.Detect()
	assert.Error(t, err)

	det = NewDetector(func() (io.Reader, error) { return strings.NewReader("ID=debian"), nil }, func(string) string { return "" }, func() domain.DesktopEnvironment { return domain.Other })
	_, err = det.Detect()
	assert.Error(t, err)
}
