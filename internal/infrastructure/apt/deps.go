package apt

import (
	"github.com/so-install/internal/core/domain"
)

// BaseDepsInstaller handles installation of essential system tools.
type BaseDepsInstaller struct {
	executor domain.Executor
}

// NewBaseDepsInstaller creates a new BaseDepsInstaller.
func NewBaseDepsInstaller(executor domain.Executor) *BaseDepsInstaller {
	return &BaseDepsInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*BaseDepsInstaller)(nil)

func (i *BaseDepsInstaller) Install() error {
	pkgs := []string{"git", "wget", "curl", "ca-certificates", "gnupg", "lsb-release"}
	args := append([]string{"install", "-y"}, pkgs...)
	_, _, err := i.executor.Execute("apt-get", args...)
	return err
}

func (i *BaseDepsInstaller) IsInstalled() (bool, error) {
	// Simple check: if git is installed, we assume base deps are mostly there
	// We could check all of them, but git is a good proxy.
	_, _, err := i.executor.Execute("which", "git")
	return err == nil, nil
}

func (i *BaseDepsInstaller) ID() domain.SoftwareID {
	return domain.BaseDeps
}
