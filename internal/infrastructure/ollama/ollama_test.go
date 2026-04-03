package ollama

import (
	"fmt"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestOllamaInstaller_IsInstalled_True(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("ollama version 0.6.0", "", nil)
	installer := &OllamaInstaller{executor: executor}

	installed, err := installer.IsInstalled()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !installed {
		t.Error("expected installed to be true")
	}
	if len(executor.Calls) != 1 || executor.Calls[0].Name != "ollama" || executor.Calls[0].Args[0] != "--version" {
		t.Errorf("unexpected executor calls: %v", executor.Calls)
	}
}

func TestOllamaInstaller_IsInstalled_False(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "not found", fmt.Errorf("exit status 127"))
	installer := &OllamaInstaller{executor: executor}

	installed, err := installer.IsInstalled()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if installed {
		t.Error("expected installed to be false")
	}
}

func TestOllamaInstaller_Install_HappyPath(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &OllamaInstaller{executor: executor}

	err := installer.Install()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(executor.Calls) != 1 {
		t.Fatalf("expected 1 executor call, got %d", len(executor.Calls))
	}
	call := executor.Calls[0]
	if call.Name != "bash" {
		t.Errorf("expected command 'bash', got %q", call.Name)
	}
	if len(call.Args) != 2 || call.Args[0] != "-c" || call.Args[1] != "curl -fsSL https://ollama.com/install.sh | sh" {
		t.Errorf("unexpected args: %v", call.Args)
	}
}

func TestOllamaInstaller_Install_Error(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "curl: network error", fmt.Errorf("exit status 1"))
	installer := &OllamaInstaller{executor: executor}

	err := installer.Install()

	if err == nil {
		t.Error("expected error, got nil")
	}
	// Should be a wrapped InstallError
	if _, ok := err.(domain.InstallError); !ok {
		if _, ok := err.(domain.AptLockError); !ok {
			t.Errorf("expected InstallError or AptLockError, got %T: %v", err, err)
		}
	}
}
