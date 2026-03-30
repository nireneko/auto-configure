package browsers

import (
	"github.com/so-install/internal/core/domain"
)

const (
	braveGPGURL     = "https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg"
	braveGPGPath    = "/usr/share/keyrings/brave-browser-archive-keyring.gpg"
	braveSourcePath = "/etc/apt/sources.list.d/brave-browser-release.list"
	braveSourceList = "deb [signed-by=/usr/share/keyrings/brave-browser-archive-keyring.gpg] https://brave-browser-apt-release.s3.brave.com/ stable main"
)

// BraveInstaller installs Brave browser from the official repository.
type BraveInstaller struct {
	executor domain.Executor
}

// NewBraveInstaller creates a new BraveInstaller.
func NewBraveInstaller(executor domain.Executor) *BraveInstaller {
	return &BraveInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*BraveInstaller)(nil)

// ID returns the SoftwareID for Brave.
func (b *BraveInstaller) ID() domain.SoftwareID { return domain.Brave }

// IsInstalled checks if brave-browser is already installed.
func (b *BraveInstaller) IsInstalled() (bool, error) {
	_, _, err := b.executor.Execute("which", "brave-browser")
	return err == nil, nil
}

// Install installs Brave from the official repository.
func (b *BraveInstaller) Install() error {
	steps := [][]string{
		{"mkdir", "-p", "/usr/share/keyrings"},
		{"wget", "-qO", braveGPGPath, braveGPGURL},
		{"sh", "-c", "echo '" + braveSourceList + "' | tee " + braveSourcePath},
		{"apt", "update"},
		{"apt", "install", "-y", "brave-browser"},
	}
	for _, step := range steps {
		_, stderr, err := b.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("brave-browser", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
