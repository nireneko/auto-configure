package antigravity

import (
	"github.com/so-install/internal/core/domain"
)

const (
	antigravityGPGURL   = "https://us-central1-apt.pkg.dev/doc/repo-signing-key.gpg"
	antigravityGPGPath  = "/etc/apt/keyrings/antigravity-repo-key.gpg"
	antigravityRepoPath = "/etc/apt/sources.list.d/antigravity.list"
	antigravityRepoEntry = `echo "deb [signed-by=/etc/apt/keyrings/antigravity-repo-key.gpg] https://us-central1-apt.pkg.dev/projects/antigravity-auto-updater-dev/ antigravity-debian main" | tee /etc/apt/sources.list.d/antigravity.list > /dev/null`
)

// AntigravityInstaller installs Google Antigravity IDE and CLI.
type AntigravityInstaller struct {
	executor domain.Executor
}

// NewAntigravityInstaller creates a new AntigravityInstaller.
func NewAntigravityInstaller(executor domain.Executor) *AntigravityInstaller {
	return &AntigravityInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*AntigravityInstaller)(nil)

// ID returns the SoftwareID for Antigravity.
func (a *AntigravityInstaller) ID() domain.SoftwareID { return domain.Antigravity }

// IsInstalled checks if antigravity (agy) is already installed.
func (a *AntigravityInstaller) IsInstalled() (bool, error) {
	_, _, err := a.executor.Execute("which", "agy")
	return err == nil, nil
}

// Install installs Google Antigravity from the official repository.
func (a *AntigravityInstaller) Install() error {
	steps := [][]string{
		{"mkdir", "-p", "/etc/apt/keyrings"},
		{"curl", "-fsSL", antigravityGPGURL, "-o", antigravityGPGPath},
		{"sh", "-c", antigravityRepoEntry},
		{"apt", "update"},
		{"apt", "install", "-y", "antigravity"},
	}
	for _, step := range steps {
		_, stderr, err := a.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("antigravity", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
