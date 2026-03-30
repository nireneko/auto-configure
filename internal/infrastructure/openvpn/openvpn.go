package openvpn

import (
	"fmt"
	"github.com/so-install/internal/core/domain"
)

// OpenVpnInstaller installs OpenVPN NetworkManager plugins.
type OpenVpnInstaller struct {
	executor domain.Executor
	osInfo   *domain.OSInfo
}

// NewOpenVpnInstaller creates a new OpenVpnInstaller.
func NewOpenVpnInstaller(executor domain.Executor, osInfo *domain.OSInfo) *OpenVpnInstaller {
	return &OpenVpnInstaller{
		executor: executor,
		osInfo:   osInfo,
	}
}

var _ domain.SoftwareInstaller = (*OpenVpnInstaller)(nil)

// ID returns the SoftwareID for OpenVPN.
func (o *OpenVpnInstaller) ID() domain.SoftwareID { return domain.OpenVpn }

// IsInstalled checks if openvpn packages are already installed.
func (o *OpenVpnInstaller) IsInstalled() (bool, error) {
	pkgs := o.getPackages()
	if len(pkgs) == 0 {
		return false, nil
	}
	// Simplified check for now: if all are installed via dpkg
	for _, pkg := range pkgs {
		_, _, err := o.executor.Execute("dpkg", "-s", pkg)
		if err != nil {
			return false, nil
		}
	}
	return true, nil
}

// Install installs OpenVPN packages and restarts NetworkManager.
func (o *OpenVpnInstaller) Install() error {
	pkgs := o.getPackages()
	if len(pkgs) == 0 {
		return fmt.Errorf("could not determine OpenVPN packages for DE: %s", o.osInfo.DesktopEnvironment)
	}

	steps := [][]string{
		{"apt-get", "update"},
		append([]string{"apt-get", "install", "-y"}, pkgs...),
		{"systemctl", "restart", "NetworkManager"},
	}

	for _, step := range steps {
		_, stderr, err := o.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("openvpn", step[0], step[1:], "", stderr, err)
		}
	}

	return nil
}

func (o *OpenVpnInstaller) getPackages() []string {
	common := []string{"network-manager-openvpn"}
	switch o.osInfo.DesktopEnvironment {
	case domain.KDE:
		return append(common, "plasma-nm")
	case domain.GNOME:
		return append(common, "network-manager-openvpn-gnome")
	default:
		return common // basic OpenVPN nm support
	}
}
