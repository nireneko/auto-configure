package homebrew

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestHomebrewInstaller_ID(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := NewHomebrewInstaller(executor)

	if installer.ID() != domain.Homebrew {
		t.Errorf("expected ID %s, got %s", domain.Homebrew, installer.ID())
	}
}

func TestHomebrewInstaller_IsInstalled(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name     string
		setup    func(path string)
		expected bool
	}{
		{
			name: "not installed",
			setup: func(path string) {
				// No brew binary
			},
			expected: false,
		},
		{
			name: "installed",
			setup: func(path string) {
				os.MkdirAll(filepath.Dir(path), 0755)
				os.WriteFile(path, []byte("fake brew"), 0755)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &mocks.MockExecutor{}
			installer := NewHomebrewInstaller(executor)
			
			brewFile := filepath.Join(tempDir, "bin/brew")
			tt.setup(brewFile)
			
			installer.brewPath = brewFile

			installed, err := installer.IsInstalled()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if installed != tt.expected {
				t.Errorf("expected installed %v, got %v", tt.expected, installed)
			}
		})
	}
}

func TestHomebrewInstaller_Install(t *testing.T) {
	tempHome := t.TempDir()
	
	// Create dummy .bashrc
	bashrcPath := filepath.Join(tempHome, ".bashrc")
	os.WriteFile(bashrcPath, []byte("# dummy bashrc\n"), 0644)

	executor := &mocks.MockExecutor{}
	installer := NewHomebrewInstaller(executor)
	installer.homeDir = tempHome

	err := installer.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify shell config
	bashrcContent, _ := os.ReadFile(bashrcPath)
	expectedLine := `eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"`
	if !strings.Contains(string(bashrcContent), expectedLine) {
		t.Errorf("expected %s in .bashrc, got %s", expectedLine, string(bashrcContent))
	}
}
