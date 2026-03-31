package osrelease

import (
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
	assert.Equal(t, "12", kvs["VERSION_ID"]) // Quotes stripped
	assert.NotContains(t, kvs, "NO_EQUALS")
}

func TestNewDefaultDetector(t *testing.T) {
	d := NewDefaultDetector()
	assert.NotNil(t, d)
}

func TestIsProcessRunning(t *testing.T) {
	// We check for a process that definitely does not exist
	assert.False(t, isProcessRunning("definitely-not-running-12345"))
}

func TestIsPackageInstalled(t *testing.T) {
	// Check for something that doesn't exist
	assert.False(t, isPackageInstalled("definitely-not-installed-12345"))
}

func TestDetectDesktopEnvironment(t *testing.T) {
	// This hits the fallback paths, likely returning Other on a CI environment
	de := detectDesktopEnvironment()
	assert.Contains(t, []domain.DesktopEnvironment{domain.GNOME, domain.KDE, domain.Other}, de)
}
