package flatpak

import (
	"fmt"

	"github.com/so-install/internal/core/domain"
)

// FlatpakInstaller installs Flatpak and its Flathub repository.
type FlatpakInstaller struct {
	executor   domain.Executor
	osDetector domain.OSDetector
}

// NewFlatpakInstaller creates a new FlatpakInstaller.
func NewFlatpakInstaller(executor domain.Executor, osDetector domain.OSDetector) *FlatpakInstaller {
	return &FlatpakInstaller{executor: executor, osDetector: osDetector}
}

var _ domain.SoftwareInstaller = (*FlatpakInstaller)(nil)

// ID returns the SoftwareID for Flatpak.
func (f *FlatpakInstaller) ID() domain.SoftwareID { return domain.Flatpak }

// IsInstalled checks if flatpak is already installed.
func (f *FlatpakInstaller) IsInstalled() (bool, error) {
	_, _, err := f.executor.Execute("flatpak", "--version")
	return err == nil, nil
}

// Install installs Flatpak, Flathub repo, and DE-specific plugins.
func (f *FlatpakInstaller) Install() error {
	steps := [][]string{
		{"apt", "update"},
		{"apt", "install", "-y", "flatpak"},
		{"flatpak", "remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"},
	}

	for _, step := range steps {
		_, stderr, err := f.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("flatpak", step[0], step[1:], "", stderr, err)
		}
	}

	// DE-specific plugins
	osInfo, err := f.osDetector.Detect()
	if err == nil && osInfo != nil {
		var plugin string
		switch osInfo.DesktopEnvironment {
		case domain.KDE:
			plugin = "plasma-discover-backend-flatpak"
		case domain.GNOME:
			plugin = "gnome-software-plugin-flatpak"
		}

		if plugin != "" {
			_, stderr, err := f.executor.Execute("apt", "install", "-y", plugin)
			if err != nil {
				return domain.WrapInstallError("flatpak", "apt", []string{"install", "-y", plugin}, "", stderr, err)
			}
		}
	}

	fmt.Println("\n[INFO] Flatpak installation completed successfully.")
	fmt.Println("[INFO] IMPORTANT: Please RESTART your session or system for changes to take effect and apps to appear in menus.")

	return nil
}
