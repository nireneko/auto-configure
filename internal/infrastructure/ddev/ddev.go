package ddev

import (
	"github.com/so-install/internal/core/domain"
)

const (
	ddevGPGURL  = "https://pkg.ddev.com/apt/gpg.key"
	ddevGPGPath = "/etc/apt/keyrings/ddev.gpg"
	ddevRepoEntry = `echo "deb [signed-by=/etc/apt/keyrings/ddev.gpg] https://pkg.ddev.com/apt/ * *" | sudo tee /etc/apt/sources.list.d/ddev.list >/dev/null`
)

// DdevInstaller installs DDEV from the official DDEV repository.
type DdevInstaller struct {
	executor domain.Executor
}

// NewDdevInstaller creates a new DdevInstaller.
func NewDdevInstaller(executor domain.Executor) *DdevInstaller {
	return &DdevInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*DdevInstaller)(nil)

// ID returns the SoftwareID for DDEV.
func (d *DdevInstaller) ID() domain.SoftwareID { return domain.Ddev }

// IsInstalled checks if ddev is already installed.
func (d *DdevInstaller) IsInstalled() (bool, error) {
	_, _, err := d.executor.Execute("ddev", "--version")
	return err == nil, nil
}

// Install installs DDEV using official steps.
func (d *DdevInstaller) Install() error {
	steps := [][]string{
		{"apt-get", "update"},
		{"apt-get", "install", "-y", "curl"},
		{"install", "-m", "0755", "-d", "/etc/apt/keyrings"},
		{"sh", "-c", "curl -fsSL " + ddevGPGURL + " | gpg --dearmor | sudo tee " + ddevGPGPath + " > /dev/null"},
		{"chmod", "a+r", ddevGPGPath},
		{"sh", "-c", ddevRepoEntry},
		{"apt-get", "update"},
		{"apt-get", "install", "-y", "ddev"},
		{"mkcert", "-install"},
	}

	for _, step := range steps {
		_, stderr, err := d.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("ddev", step[0], step[1:], "", stderr, err)
		}
	}

	return nil
}
