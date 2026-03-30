package flatpak

import (
	"github.com/so-install/internal/core/domain"
)

// FlatpakAppInstaller installs applications from Flathub.
type FlatpakAppInstaller struct {
	executor   domain.Executor
	appID      string
	softwareID domain.SoftwareID
}

// NewFlatpakAppInstaller creates a new FlatpakAppInstaller.
func NewFlatpakAppInstaller(executor domain.Executor, appID string, id domain.SoftwareID) *FlatpakAppInstaller {
	return &FlatpakAppInstaller{
		executor:   executor,
		appID:      appID,
		softwareID: id,
	}
}

var _ domain.SoftwareInstaller = (*FlatpakAppInstaller)(nil)

// ID returns the SoftwareID.
func (f *FlatpakAppInstaller) ID() domain.SoftwareID { return f.softwareID }

// IsInstalled checks if the application is already installed using flatpak info.
func (f *FlatpakAppInstaller) IsInstalled() (bool, error) {
	_, _, err := f.executor.Execute("flatpak", "info", f.appID)
	return err == nil, nil
}

// Install installs the application from Flathub.
func (f *FlatpakAppInstaller) Install() error {
	_, stderr, err := f.executor.Execute("flatpak", "install", "-y", "flathub", f.appID)
	if err != nil {
		return domain.WrapInstallError(string(f.softwareID), "flatpak", []string{"install", "-y", "flathub", f.appID}, "", stderr, err)
	}
	return nil
}
