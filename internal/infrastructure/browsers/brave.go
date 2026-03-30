package browsers

import (
	"strings"

	"github.com/so-install/internal/core/domain"
)

const (
	braveGPGURL     = "https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg"
	braveSourcesURL = "https://brave-browser-apt-release.s3.brave.com/brave-browser.sources"
	braveGPGPath    = "/usr/share/keyrings/brave-browser-archive-keyring.gpg"
	braveSourcePath = "/etc/apt/sources.list.d/brave-browser-release.sources"
)

// BraveInstaller installs Brave browser from the official repository.
type BraveInstaller struct {
	executor domain.Executor
}

// NewBraveInstaller creates a new BraveInstaller.
func NewBraveInstaller(executor domain.Executor) *BraveInstaller {
	return &BraveInstaller{executor: executor}
}

var _ domain.BrowserInstaller = (*BraveInstaller)(nil)

// ID returns the BrowserID for Brave.
func (b *BraveInstaller) ID() domain.BrowserID { return domain.Brave }

// IsInstalled checks if brave-browser is already installed.
func (b *BraveInstaller) IsInstalled() (bool, error) {
	_, _, err := b.executor.Execute("which", "brave-browser")
	return err == nil, nil
}

// Install installs Brave from the official repository.
func (b *BraveInstaller) Install() error {
	steps := [][]string{
		{"curl", "-fsSLo", braveGPGPath, braveGPGURL},
		{"curl", "-fsSLo", braveSourcePath, braveSourcesURL},
		{"apt", "update"},
		{"apt", "install", "-y", "brave-browser"},
	}
	for _, step := range steps {
		_, stderr, err := b.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return wrapInstallError("brave-browser", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}

// wrapInstallError wraps a shell error as InstallError or AptLockError.
func wrapInstallError(browser, cmd string, args []string, stdout, stderr string, err error) error {
	base := domain.InstallError{
		Browser: browser,
		Command: cmd,
		Args:    args,
		Stdout:  stdout,
		Stderr:  stderr,
	}
	if strings.Contains(stderr, "Could not get lock") {
		return domain.AptLockError{InstallError: base}
	}
	return base
}
