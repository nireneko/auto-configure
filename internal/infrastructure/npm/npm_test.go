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
		installer := &NpmInstaller{
			executor:    executor,
			packageName: "test-pkg",
			binaryName:  "test-bin",
			softwareID:  "test",
			homeDir:     "/nonexistent", // Force no nvm sourcing
			userName:    "",             // No sudo -u
		}

		installed, err := installer.IsInstalled()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !installed {
			t.Error("Expected installed to be true")
		}
		
		// If no SUDO_USER, it uses bash -c "test-bin --version"
		if executor.Calls[0].Name != "bash" || executor.Calls[0].Args[1] != "test-bin --version" {
			t.Errorf("Expected call to bash -c 'test-bin --version', got %s %v", executor.Calls[0].Name, executor.Calls[0].Args)
		}
	})

	t.Run("should return false if command fails", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "not found", fmt.Errorf("error"))
		installer := &NpmInstaller{
			executor:    executor,
			packageName: "test-pkg",
			binaryName:  "test-bin",
			softwareID:  "test",
			homeDir:     "/nonexistent",
			userName:    "",
		}

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
		installer := &NpmInstaller{
			executor:    executor,
			packageName: "@org/pkg",
			binaryName:  "pkg-bin",
			softwareID:  "test",
			homeDir:     "/nonexistent",
			userName:    "",
		}
		
		err := installer.Install()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(executor.Calls) != 1 {
			t.Fatalf("Expected 1 call, got %d", len(executor.Calls))
		}

		call := executor.Calls[0]
		if call.Name != "bash" || call.Args[1] != "npm install -g @org/pkg" {
			t.Errorf("Unexpected command: %s %v", call.Name, call.Args)
		}
	})

	t.Run("should source nvm and use sudo if nvm.sh exists and userName is provided", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		
		// Create a temporary nvm.sh to trigger the sourcing logic
		tmpDir := t.TempDir()
		nvmDir := filepath.Join(tmpDir, ".nvm")
		os.MkdirAll(nvmDir, 0755)
		nvmScript := filepath.Join(nvmDir, "nvm.sh")
		os.WriteFile(nvmScript, []byte("#!/bin/bash"), 0644)
		
		installer := &NpmInstaller{
			executor:    executor,
			packageName: "test-pkg",
			binaryName:  "test-bin",
			softwareID:  "test",
			homeDir:     tmpDir,
			userName:    "testuser",
		}

		err := installer.Install()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(executor.Calls) != 1 {
			t.Fatalf("Expected 1 call, got %d", len(executor.Calls))
		}

		call := executor.Calls[0]
		if call.Name != "sudo" {
			t.Errorf("Expected sudo, got %s", call.Name)
		}
		if call.Args[1] != "testuser" {
			t.Errorf("Expected -u testuser, got -u %s", call.Args[1])
		}
		
		expectedCmd := fmt.Sprintf("source %s && npm install -g test-pkg", nvmScript)
		if call.Args[4] != expectedCmd {
			t.Errorf("Expected command %q, got %q", expectedCmd, call.Args[4])
		}
	})
}
