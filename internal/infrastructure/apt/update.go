package apt

import (
	"github.com/so-install/internal/core/domain"
)

// AptUpdateInstaller handles system updates and upgrades.
type AptUpdateInstaller struct {
	executor domain.Executor
}

// NewAptUpdateInstaller creates a new AptUpdateInstaller.
func NewAptUpdateInstaller(executor domain.Executor) *AptUpdateInstaller {
	return &AptUpdateInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*AptUpdateInstaller)(nil)

func (i *AptUpdateInstaller) Install() error {
	// 1. apt-get update
	_, _, err := i.executor.Execute("apt-get", "update")
	if err != nil {
		return err
	}

	// 2. apt-get upgrade -y
	_, _, err = i.executor.Execute("sh", "-c", "DEBIAN_FRONTEND=noninteractive apt-get upgrade -y")
	return err
}

func (i *AptUpdateInstaller) IsInstalled() (bool, error) {
	// Always returns false as we want it to run at the start of each execution
	return false, nil
}

func (i *AptUpdateInstaller) ID() domain.SoftwareID {
	return domain.SystemUpdate
}
