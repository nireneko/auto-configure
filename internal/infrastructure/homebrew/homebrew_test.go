package homebrew

import (
	"fmt"
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
	t.Run("should install successfully", func(t *testing.T) {
		tempHome := t.TempDir()
		
		// Create dummy .bashrc
		bashrcPath := filepath.Join(tempHome, ".bashrc")
		os.WriteFile(bashrcPath, []byte("# dummy bashrc\n"), 0644)

		executor := &mocks.MockExecutor{}
		installer := &HomebrewInstaller{
			executor: executor,
			homeDir:  tempHome,
			brewPath: "/tmp/fakebrew",
			userName: "testuser",
		}

		err := installer.Install()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(executor.Calls) < 2 {
			t.Fatalf("expected at least 2 calls, got %d", len(executor.Calls))
		}

		if executor.Calls[0].Name != "apt" {
			t.Errorf("expected call to apt, got %s", executor.Calls[0].Name)
		}

		if executor.Calls[1].Name != "sudo" || executor.Calls[1].Args[1] != "testuser" {
			t.Errorf("expected sudo -u testuser for brew install, got %s %v", executor.Calls[1].Name, executor.Calls[1].Args)
		}

		// Verify shell config
		bashrcContent, _ := os.ReadFile(bashrcPath)
		expectedLine := `eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"`
		if !strings.Contains(string(bashrcContent), expectedLine) {
			t.Errorf("expected %s in .bashrc, got %s", expectedLine, string(bashrcContent))
		}
	})

	t.Run("should return error if apt install fails", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "apt failed", fmt.Errorf("error apt"))

		installer := &HomebrewInstaller{
			executor: executor,
			homeDir:  t.TempDir(),
			brewPath: "/tmp/fakebrew",
			userName: "testuser",
		}

		err := installer.Install()
		if err == nil {
			t.Error("expected error from apt install, got nil")
		}
	})

	t.Run("should return error if brew install script fails", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		// apt succeeds
		executor.AddResponse("", "", nil)
		// brew script fails
		executor.AddResponse("", "brew script failed", fmt.Errorf("error script"))

		installer := &HomebrewInstaller{
			executor: executor,
			homeDir:  t.TempDir(),
			brewPath: "/tmp/fakebrew",
			userName: "root", // tests fallback to bash if root
		}

		err := installer.Install()
		if err == nil {
			t.Error("expected error from brew install script, got nil")
		}
	})

	t.Run("should skip modify if brew already in bashrc", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		tempHome := t.TempDir()
		bashrcPath := filepath.Join(tempHome, ".bashrc")
		expectedLine := `eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"`
		os.WriteFile(bashrcPath, []byte(expectedLine), 0644)

		installer := &HomebrewInstaller{
			executor: executor,
			homeDir:  tempHome,
			brewPath: "/tmp/fakebrew",
			userName: "testuser",
		}
		
		installer.Install()
	})
}
