package openvpn_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/openvpn"
	"github.com/so-install/pkg/mocks"
)

func TestOpenVpnInstaller_Install_KDE(t *testing.T) {
	executor := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{DesktopEnvironment: domain.KDE}
	installer := openvpn.NewOpenVpnInstaller(executor, osInfo)

	err := installer.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify packages for KDE
	foundPkgs := false
	foundRestart := false
	for _, call := range executor.Calls {
		if call.Name == "apt-get" && len(call.Args) >= 4 && call.Args[0] == "install" && call.Args[1] == "-y" && call.Args[2] == "network-manager-openvpn" && call.Args[3] == "plasma-nm" {
			foundPkgs = true
		}
		if call.Name == "systemctl" && len(call.Args) == 2 && call.Args[0] == "restart" && call.Args[1] == "NetworkManager" {
			foundRestart = true
		}
	}

	if !foundPkgs {
		t.Errorf("expected KDE packages not found in execution log")
	}
	if !foundRestart {
		t.Errorf("expected NetworkManager restart not found in execution log")
	}
}

func TestOpenVpnInstaller_Install_GNOME(t *testing.T) {
	executor := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{DesktopEnvironment: domain.GNOME}
	installer := openvpn.NewOpenVpnInstaller(executor, osInfo)

	err := installer.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify packages for GNOME
	foundPkgs := false
	foundRestart := false
	for _, call := range executor.Calls {
		if call.Name == "apt-get" && len(call.Args) >= 4 && call.Args[0] == "install" && call.Args[1] == "-y" && call.Args[2] == "network-manager-openvpn" && call.Args[3] == "network-manager-openvpn-gnome" {
			foundPkgs = true
		}
		if call.Name == "systemctl" && len(call.Args) == 2 && call.Args[0] == "restart" && call.Args[1] == "NetworkManager" {
			foundRestart = true
		}
	}

	if !foundPkgs {
		t.Errorf("expected GNOME packages not found in execution log")
	}
	if !foundRestart {
		t.Errorf("expected NetworkManager restart not found in execution log")
	}
}

func TestOpenVpnInstaller_Install_Error(t *testing.T) {
	executor := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{DesktopEnvironment: domain.KDE}
	installer := openvpn.NewOpenVpnInstaller(executor, osInfo)

	executor.AddResponse("", "apt update failed", domain.WrapInstallError("openvpn", "apt-get", []string{"update"}, "", "apt update failed", nil))

	err := installer.Install()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestOpenVpnInstaller_Install_NoDE(t *testing.T) {
	executor := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{DesktopEnvironment: ""}
	installer := openvpn.NewOpenVpnInstaller(executor, osInfo)

	err := installer.Install()
	if err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}
}

func TestOpenVpnInstaller_IsInstalled(t *testing.T) {
	executor := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{DesktopEnvironment: domain.KDE}
	installer := openvpn.NewOpenVpnInstaller(executor, osInfo)

	executor.AddResponse("", "", nil) // dpkg -s network-manager-openvpn
	executor.AddResponse("", "", nil) // dpkg -s plasma-nm

	installed, err := installer.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !installed {
		t.Fatalf("expected to be installed")
	}
}

func TestOpenVpnInstaller_IsNotInstalled(t *testing.T) {
	executor := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{DesktopEnvironment: domain.KDE}
	installer := openvpn.NewOpenVpnInstaller(executor, osInfo)

	executor.AddResponse("", "", domain.WrapInstallError("openvpn", "dpkg", []string{"-s", "network-manager-openvpn"}, "", "not installed", nil)) 

	installed, err := installer.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if installed {
		t.Fatalf("expected to not be installed")
	}
}
