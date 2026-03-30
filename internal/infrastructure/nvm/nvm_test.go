package nvm

import (
	"fmt"
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
}

func TestNvmInstaller_Install(t *testing.T) {
	t.Run("should execute install commands", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewNvmInstaller(executor)
		installer.homeDir = "/tmp/fakehome"

		err := installer.Install()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		if len(executor.Calls) != 2 {
			t.Fatalf("Expected 2 calls, got %d", len(executor.Calls))
		}

		expectedInstallCmd := fmt.Sprintf("curl -o- %s | bash", nvmInstallURL)
		if executor.Calls[0].Name != "bash" || executor.Calls[0].Args[1] != expectedInstallCmd {
			t.Errorf("Expected command 'bash -c %s' not found", expectedInstallCmd)
		}

		expectedNodeCmd := "source /tmp/fakehome/.nvm/nvm.sh && nvm install --lts"
		if executor.Calls[1].Name != "bash" || executor.Calls[1].Args[1] != expectedNodeCmd {
			t.Errorf("Expected command 'bash -c %s' not found", expectedNodeCmd)
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
}
