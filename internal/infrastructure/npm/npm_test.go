package npm

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
)

func TestNpmInstaller_ID(t *testing.T) {
	executor := &mocks.MockExecutor{}
	id := domain.SoftwareID("test-pkg")
	installer := NewNpmInstaller(executor, "test-pkg", "test-bin", id)
	if installer.ID() != id {
		t.Errorf("Expected ID %s, got %s", id, installer.ID())
	}
}

func TestNpmInstaller_IsInstalled(t *testing.T) {
	t.Run("should return true if command succeeds", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewNpmInstaller(executor, "test-pkg", "test-bin", "test")
		installer.homeDir = "/nonexistent" // Force no nvm sourcing

		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !installed {
			t.Error("Expected installed to be true")
		}
		
		if executor.Calls[0].Name != "test-bin" {
			t.Errorf("Expected call to test-bin, got %s", executor.Calls[0].Name)
		}
	})

	t.Run("should return false if command fails", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "not found", fmt.Errorf("error"))
		installer := NewNpmInstaller(executor, "test-pkg", "test-bin", "test")
		installer.homeDir = "/nonexistent"

		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if installed {
			t.Error("Expected installed to be false")
		}
	})
}

func TestNpmInstaller_Install(t *testing.T) {
	t.Run("should execute npm install -g", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewNpmInstaller(executor, "@org/pkg", "pkg-bin", "test")
		installer.homeDir = "/nonexistent"

		err := installer.Install()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(executor.Calls) != 1 {
			t.Fatalf("Expected 1 call, got %d", len(executor.Calls))
		}

		call := executor.Calls[0]
		if call.Name != "npm" || call.Args[0] != "install" || call.Args[1] != "-g" || call.Args[2] != "@org/pkg" {
			t.Errorf("Unexpected command: %s %v", call.Name, call.Args)
		}
	})

	t.Run("should source nvm if nvm.sh exists", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewNpmInstaller(executor, "test-pkg", "test-bin", "test")
		
		// Create a temporary nvm.sh to trigger the sourcing logic
		tmpDir := t.TempDir()
		nvmDir := filepath.Join(tmpDir, ".nvm")
		os.MkdirAll(nvmDir, 0755)
		nvmScript := filepath.Join(nvmDir, "nvm.sh")
		os.WriteFile(nvmScript, []byte("#!/bin/bash"), 0644)
		
		installer.homeDir = tmpDir

		err := installer.Install()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(executor.Calls) != 1 {
			t.Fatalf("Expected 1 call, got %d", len(executor.Calls))
		}

		call := executor.Calls[0]
		if call.Name != "bash" {
			t.Errorf("Expected bash, got %s", call.Name)
		}
		expectedCmd := fmt.Sprintf("source %s && npm install -g test-pkg", nvmScript)
		if call.Args[1] != expectedCmd {
			t.Errorf("Expected command %q, got %q", expectedCmd, call.Args[1])
		}
	})
}
