package opencode

import (
	"fmt"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestOpenCodeInstaller_IsInstalled_True(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("0.1.0", "", nil)
	installer := &OpenCodeInstaller{executor: executor, userName: "alice"}

	installed, err := installer.IsInstalled()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !installed {
		t.Error("expected installed to be true")
	}
	if len(executor.Calls) != 1 || executor.Calls[0].Name != "opencode" || executor.Calls[0].Args[0] != "--version" {
		t.Errorf("unexpected executor calls: %v", executor.Calls)
	}
}

func TestOpenCodeInstaller_IsInstalled_False(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "not found", fmt.Errorf("exit status 127"))
	installer := &OpenCodeInstaller{executor: executor, userName: "alice"}

	installed, err := installer.IsInstalled()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if installed {
		t.Error("expected installed to be false")
	}
}

func TestOpenCodeInstaller_Install_NonRootUser(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &OpenCodeInstaller{executor: executor, userName: "alice"}

	err := installer.Install()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(executor.Calls) != 1 {
		t.Fatalf("expected 1 executor call, got %d", len(executor.Calls))
	}
	call := executor.Calls[0]
	if call.Name != "sudo" {
		t.Errorf("expected command 'sudo', got %q", call.Name)
	}
	expectedArgs := []string{"-u", "alice", "bash", "-c", "curl -fsSL https://opencode.ai/install | bash"}
	if len(call.Args) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs), len(call.Args), call.Args)
	}
	for i, arg := range expectedArgs {
		if call.Args[i] != arg {
			t.Errorf("arg[%d]: expected %q, got %q", i, arg, call.Args[i])
		}
	}
}

func TestOpenCodeInstaller_Install_EmptyUser(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &OpenCodeInstaller{executor: executor, userName: ""}

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
	expectedArgs := []string{"-c", "curl -fsSL https://opencode.ai/install | bash"}
	if len(call.Args) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs), len(call.Args), call.Args)
	}
	for i, arg := range expectedArgs {
		if call.Args[i] != arg {
			t.Errorf("arg[%d]: expected %q, got %q", i, arg, call.Args[i])
		}
	}
}

func TestOpenCodeInstaller_Install_RootUser(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &OpenCodeInstaller{executor: executor, userName: "root"}

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
}

func TestOpenCodeInstaller_Install_Error(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "network error", fmt.Errorf("exit status 1"))
	installer := &OpenCodeInstaller{executor: executor, userName: "alice"}

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
