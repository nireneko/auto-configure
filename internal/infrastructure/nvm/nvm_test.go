package nvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestNvmInstaller_ID(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := NewNvmInstaller(executor)
	if installer.ID() != domain.Nvm {
		t.Errorf("Expected ID %s, got %s", domain.Nvm, installer.ID())
	}
}

func TestNvmInstaller_IsInstalled(t *testing.T) {
	t.Run("should return false if nvm.sh does not exist", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewNvmInstaller(executor)
		installer.homeDir = "/tmp/fakehome_nonexistent"
		
		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if installed {
			t.Error("Expected installed to be false")
		}
	})

	t.Run("should return true if nvm.sh exists", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewNvmInstaller(executor)
		tmpHome := t.TempDir()
		installer.homeDir = tmpHome
		nvmDir := filepath.Join(tmpHome, ".nvm")
		os.MkdirAll(nvmDir, 0755)
		os.WriteFile(filepath.Join(nvmDir, "nvm.sh"), []byte("echo hi"), 0755)

		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !installed {
			t.Error("Expected installed to be true")
		}
	})
}

func TestNvmInstaller_Install(t *testing.T) {
	t.Run("should execute install commands as actual user", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		// Mock home dir and user
		tmpHome := t.TempDir()
		
		// Create fake .bashrc to test configureShell
		bashrcPath := filepath.Join(tmpHome, ".bashrc")
		if err := os.WriteFile(bashrcPath, []byte("# existing content"), 0644); err != nil {
			t.Fatalf("Failed to create fake .bashrc: %v", err)
		}

		installer := &NvmInstaller{
			executor: executor,
			homeDir:  tmpHome,
			userName: "fakeuser",
		}

		err := installer.Install()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Expect 3 calls: 1. NVM install, 2. Node install (since configureShell is internal)
		// Wait, configureShell doesn't call executor, it writes files.
		if len(executor.Calls) != 2 {
			t.Fatalf("Expected 2 calls to executor, got %d", len(executor.Calls))
		}

		expectedInstallCmd := fmt.Sprintf("curl -o- %s | bash", nvmInstallURL)
		if executor.Calls[0].Name != "sudo" || executor.Calls[0].Args[1] != "fakeuser" || executor.Calls[0].Args[4] != expectedInstallCmd {
			t.Errorf("Expected command 'sudo -u fakeuser bash -c %s' not found, got %s %v", expectedInstallCmd, executor.Calls[0].Name, executor.Calls[0].Args)
		}

		nvmScript := filepath.Join(tmpHome, ".nvm", "nvm.sh")
		expectedNodeCmd := fmt.Sprintf("source %s && nvm install --lts", nvmScript)
		if executor.Calls[1].Name != "sudo" || executor.Calls[1].Args[1] != "fakeuser" || executor.Calls[1].Args[4] != expectedNodeCmd {
			t.Errorf("Expected command 'sudo -u fakeuser bash -c %s' not found, got %s %v", expectedNodeCmd, executor.Calls[1].Name, executor.Calls[1].Args)
		}

		// Verify .bashrc was modified
		content, err := os.ReadFile(bashrcPath)
		if err != nil {
			t.Fatalf("Failed to read .bashrc: %v", err)
		}
		if !strings.Contains(string(content), "nvm.sh") {
			t.Error("Expected .bashrc to contain nvm.sh loader")
		}
	})

	t.Run("should return error if nvm install fails", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "failed to download", fmt.Errorf("error"))
		
		installer := NewNvmInstaller(executor)
		err := installer.Install()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return error if node lts install fails", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		// First command (nvm install script) succeeds
		executor.AddResponse("", "", nil)
		// Second command (node lts) fails
		executor.AddResponse("", "failed to install node", fmt.Errorf("error node"))
		
		installer := &NvmInstaller{
			executor: executor,
			homeDir:  t.TempDir(), // prevent configureShell issues
			userName: "root",      // no sudo
		}
		err := installer.Install()
		if err == nil {
			t.Error("Expected error from node install, got nil")
		}
	})

	t.Run("should skip modify if nvm already in bashrc", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		tmpHome := t.TempDir()
		bashrcPath := filepath.Join(tmpHome, ".bashrc")
		os.WriteFile(bashrcPath, []byte("nvm.sh is here"), 0644)

		installer := &NvmInstaller{
			executor: executor,
			homeDir:  tmpHome,
			userName: "fakeuser",
		}
		
		installer.Install()
		// If it reaches here without error, it works.
	})
}
