package nvm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/so-install/internal/core/domain"
)

const (
	nvmInstallURL = "https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh"
)

// NvmInstaller installs NVM and Node.js LTS.
type NvmInstaller struct {
	executor domain.Executor
	homeDir  string
}

// NewNvmInstaller creates a new NvmInstaller.
func NewNvmInstaller(executor domain.Executor) *NvmInstaller {
	home, _ := os.UserHomeDir()
	return &NvmInstaller{
		executor: executor,
		homeDir:  home,
	}
}

var _ domain.SoftwareInstaller = (*NvmInstaller)(nil)

// ID returns the SoftwareID for NVM.
func (n *NvmInstaller) ID() domain.SoftwareID { return domain.Nvm }

// IsInstalled checks if NVM is already installed by looking for nvm.sh.
func (n *NvmInstaller) IsInstalled() (bool, error) {
	nvmScript := filepath.Join(n.homeDir, ".nvm", "nvm.sh")
	_, err := os.Stat(nvmScript)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Install installs NVM and Node.js LTS.
func (n *NvmInstaller) Install() error {
	// 1. Install NVM using official script
	installCmd := fmt.Sprintf("curl -o- %s | bash", nvmInstallURL)
	_, stderr, err := n.executor.Execute("bash", "-c", installCmd)
	if err != nil {
		return domain.WrapInstallError("nvm", "bash", []string{"-c", installCmd}, "", stderr, err)
	}

	// 2. Install Node.js LTS using NVM
	nvmScript := filepath.Join(n.homeDir, ".nvm", "nvm.sh")
	nodeCmd := fmt.Sprintf("source %s && nvm install --lts", nvmScript)
	_, stderr, err = n.executor.Execute("bash", "-c", nodeCmd)
	if err != nil {
		return domain.WrapInstallError("nvm", "bash", []string{"-c", nodeCmd}, "", stderr, err)
	}

	return nil
}
