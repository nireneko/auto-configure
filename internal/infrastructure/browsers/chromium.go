package browsers

import "github.com/so-install/internal/core/domain"

// ChromiumInstaller installs Chromium from the Debian repository.
type ChromiumInstaller struct {
	executor domain.Executor
}

// NewChromiumInstaller creates a new ChromiumInstaller.
func NewChromiumInstaller(executor domain.Executor) *ChromiumInstaller {
	return &ChromiumInstaller{executor: executor}
}

var _ domain.BrowserInstaller = (*ChromiumInstaller)(nil)

// ID returns the BrowserID for Chromium.
func (c *ChromiumInstaller) ID() domain.BrowserID { return domain.Chromium }

// IsInstalled checks if chromium is already installed.
func (c *ChromiumInstaller) IsInstalled() (bool, error) {
	_, _, err := c.executor.Execute("which", "chromium")
	return err == nil, nil
}

// Install installs Chromium from the Debian repository.
func (c *ChromiumInstaller) Install() error {
	steps := [][]string{
		{"apt", "update"},
		{"apt", "install", "-y", "chromium"},
	}
	for _, step := range steps {
		_, stderr, err := c.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return wrapInstallError("chromium", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
