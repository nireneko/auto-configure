package apt_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/apt"
	"github.com/so-install/pkg/mocks"
)

func TestBaseDepsInstaller_Install(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewBaseDepsInstaller(mockExecutor)

	// We expect 1 call to apt-get install -y git wget curl ca-certificates gnupg lsb-release
	mockExecutor.AddResponse("", "", nil)

	err := installer.Install()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockExecutor.Calls) != 1 {
		t.Fatalf("Expected 1 call to executor, got %d", len(mockExecutor.Calls))
	}

	call := mockExecutor.Calls[0]
	if call.Name != "apt-get" {
		t.Errorf("Call = %s, want apt-get", call.Name)
	}

	expectedArgs := []string{"install", "-y", "git", "wget", "curl", "ca-certificates", "gnupg", "lsb-release"}
	for _, expected := range expectedArgs {
		found := false
		for _, arg := range call.Args {
			if arg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Call args = %v, missing expected arg %s", call.Args, expected)
		}
	}
}

func TestBaseDepsInstaller_IsInstalled(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewBaseDepsInstaller(mockExecutor)

	// We expect 1 call to check if git is installed as a proxy for all base deps
	mockExecutor.AddResponse("", "", nil)

	installed, err := installer.IsInstalled()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !installed {
		t.Error("Expected installed to be true when which git returns no error")
	}

	if mockExecutor.Calls[0].Name != "which" || mockExecutor.Calls[0].Args[0] != "git" {
		t.Errorf("Call = %s %v, want which git", mockExecutor.Calls[0].Name, mockExecutor.Calls[0].Args)
	}
}

func TestBaseDepsInstaller_ID(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewBaseDepsInstaller(mockExecutor)

	if installer.ID() != domain.BaseDeps {
		t.Errorf("ID = %v, want %v", installer.ID(), domain.BaseDeps)
	}
}
