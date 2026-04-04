package cursor

import (
	"github.com/so-install/internal/core/domain"
)

const (
	cursorURL = "https://downloader.cursor.sh/linux/debian/amd64"
	cursorDeb = "/tmp/cursor.deb"
)

// CursorInstaller installs Cursor IDE from the official download.
type CursorInstaller struct {
	executor domain.Executor
}

// NewCursorInstaller creates a new CursorInstaller.
func NewCursorInstaller(executor domain.Executor) *CursorInstaller {
	return &CursorInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*CursorInstaller)(nil)

// ID returns the SoftwareID for Cursor.
func (c *CursorInstaller) ID() domain.SoftwareID { return domain.Cursor }

// IsInstalled checks if cursor is already installed.
func (c *CursorInstaller) IsInstalled() (bool, error) {
	_, _, err := c.executor.Execute("which", "cursor")
	return err == nil, nil
}

// Install downloads and installs Cursor IDE.
func (c *CursorInstaller) Install() error {
	steps := [][]string{
		{"wget", cursorURL, "-O", cursorDeb},
		{"apt", "install", "-y", cursorDeb},
	}
	for _, step := range steps {
		_, stderr, err := c.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("cursor", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
