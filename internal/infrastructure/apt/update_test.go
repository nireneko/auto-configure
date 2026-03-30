package apt_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/apt"
	"github.com/so-install/pkg/mocks"
)

func TestAptUpdateInstaller_Install(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewAptUpdateInstaller(mockExecutor)

	// We expect 2 calls: apt-get update and apt-get upgrade -y
	mockExecutor.AddResponse("", "", nil) // update ok
	mockExecutor.AddResponse("", "", nil) // upgrade ok

	err := installer.Install()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockExecutor.Calls) != 2 {
		t.Fatalf("Expected 2 calls to executor, got %d", len(mockExecutor.Calls))
	}

	// 1. apt-get update
	if mockExecutor.Calls[0].Name != "apt-get" || mockExecutor.Calls[0].Args[0] != "update" {
		t.Errorf("First call = %s %v, want apt-get update", mockExecutor.Calls[0].Name, mockExecutor.Calls[0].Args)
	}

	// 2. apt-get upgrade -y
	upgradeCall := mockExecutor.Calls[1]
	if upgradeCall.Name != "sh" {
		t.Errorf("Second call = %s, want sh", upgradeCall.Name)
	}
	
	foundUpgrade := false
	foundY := false
	foundFrontend := false
	for _, arg := range upgradeCall.Args {
		if arg == "DEBIAN_FRONTEND=noninteractive apt-get upgrade -y" {
			foundUpgrade = true
			foundY = true
			foundFrontend = true
		}
	}
	if !foundUpgrade || !foundY || !foundFrontend {
		t.Errorf("Second call args = %v, want to include upgrade, -y and DEBIAN_FRONTEND", upgradeCall.Args)
	}
}

func TestAptUpdateInstaller_Install_Error(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewAptUpdateInstaller(mockExecutor)

	mockExecutor.AddResponse("", "lock error", errors.New("lock failed"))

	err := installer.Install()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "lock failed" {
		t.Errorf("Error = %v, want 'lock failed'", err)
	}
}

func TestAptUpdateInstaller_IsInstalled(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewAptUpdateInstaller(mockExecutor)

	installed, err := installer.IsInstalled()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if installed {
		t.Error("AptUpdateInstaller should always return false for IsInstalled")
	}
}

func TestAptUpdateInstaller_ID(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{}
	installer := apt.NewAptUpdateInstaller(mockExecutor)

	if installer.ID() != domain.SystemUpdate {
		t.Errorf("ID = %v, want %v", installer.ID(), domain.SystemUpdate)
	}
}
