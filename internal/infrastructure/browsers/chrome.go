package browsers

import (
	"github.com/so-install/internal/core/domain"
)

const (
	chromeURL = "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb"
	chromeDeb = "/tmp/google-chrome-stable_current_amd64.deb"
)

// ChromeInstaller installs Google Chrome from the official download.
type ChromeInstaller struct {
	executor domain.Executor
}

// NewChromeInstaller creates a new ChromeInstaller.
func NewChromeInstaller(executor domain.Executor) *ChromeInstaller {
	return &ChromeInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*ChromeInstaller)(nil)

// ID returns the SoftwareID for Chrome.
func (c *ChromeInstaller) ID() domain.SoftwareID { return domain.Chrome }

// IsInstalled checks if google-chrome-stable is already installed.
func (c *ChromeInstaller) IsInstalled() (bool, error) {
	_, _, err := c.executor.Execute("which", "google-chrome-stable")
	return err == nil, nil
}

// Install downloads and installs Google Chrome.
func (c *ChromeInstaller) Install() error {
	steps := [][]string{
		{"wget", chromeURL, "-P", "/tmp/"},
		{"apt", "install", "-y", chromeDeb},
	}
	for _, step := range steps {
		_, stderr, err := c.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("google-chrome", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
