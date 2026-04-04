package homebrew

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestHomebrewInstaller_ID(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := NewHomebrewInstaller(executor)
	if installer.ID() != domain.Homebrew { t.Errorf("expected ID %s, got %s", domain.Homebrew, installer.ID()) }
}

func TestHomebrewInstaller_IsInstalled(t *testing.T) {
	tempDir := t.TempDir()
	brewFile := filepath.Join(tempDir, "brew")
	os.WriteFile(brewFile, []byte("fake"), 0755)
	installer := &HomebrewInstaller{brewPath: brewFile}
	installed, _ := installer.IsInstalled()
	if !installed { t.Error("expected true") }

	installer.brewPath = "/non/existent/path"
	installed, _ = installer.IsInstalled()
	if installed { t.Error("expected false") }
}

func TestHomebrewInstaller_Install_Full(t *testing.T) {
	tempHome := t.TempDir()
	bashrc := filepath.Join(tempHome, ".bashrc")
	zshrc := filepath.Join(tempHome, ".zshrc")
	os.WriteFile(bashrc, []byte(""), 0644)
	os.WriteFile(zshrc, []byte(""), 0644)

	executor := &mocks.MockExecutor{}
	installer := &HomebrewInstaller{
		executor: executor,
		homeDir:  tempHome,
		brewPath: "/tmp/fakebrew",
		userName: "testuser",
	}

	err := installer.Install()
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if len(executor.Calls) < 2 { t.Fatalf("expected 2 calls, got %d", len(executor.Calls)) }

	// Coverage for strings.Contains in configureShell
	installer.Install()
}

func TestHomebrewInstaller_Install_Errors(t *testing.T) {
	// apt error
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "apt fail", fmt.Errorf("err"))
	installer := &HomebrewInstaller{executor: executor}
	err := installer.Install()
	if err == nil { t.Error("expected error") }

	// script error
	executor = &mocks.MockExecutor{}
	executor.AddResponse("", "", nil) // apt ok
	executor.AddResponse("", "script fail", fmt.Errorf("err"))
	installer = &HomebrewInstaller{executor: executor, userName: "test"}
	err = installer.Install()
	if err == nil { t.Error("expected error") }
}
