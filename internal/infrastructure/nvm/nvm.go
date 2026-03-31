package nvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/so-install/internal/core/domain"
)

const (
	nvmInstallURL = "https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh"
)

// NvmInstaller installs NVM and Node.js LTS.
type NvmInstaller struct {
	executor domain.Executor
	homeDir  string
	userName string
}

// NewNvmInstaller creates a new NvmInstaller.
func NewNvmInstaller(executor domain.Executor) *NvmInstaller {
	return &NvmInstaller{
		executor: executor,
		homeDir:  domain.GetActualHome(),
		userName: domain.GetActualUser(),
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
	// 1. Install NVM using official script AS THE ACTUAL USER
	installCmd := fmt.Sprintf("curl -o- %s | bash", nvmInstallURL)
	var name string
	var args []string
	if n.userName != "" && n.userName != "root" {
		name = "sudo"
		args = []string{"-u", n.userName, "bash", "-c", installCmd}
	} else {
		name = "bash"
		args = []string{"-c", installCmd}
	}
	
	_, stderr, err := n.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError("nvm", name, args, "", stderr, err)
	}

	// 2. Configure shell environment explicitly to be sure
	if err := n.configureShell(); err != nil {
		return fmt.Errorf("failed to configure nvm shell environment: %w", err)
	}

	// 3. Install Node.js LTS using NVM
	nvmScript := filepath.Join(n.homeDir, ".nvm", "nvm.sh")
	nodeCmd := fmt.Sprintf("source %s && nvm install --lts", nvmScript)
	if n.userName != "" && n.userName != "root" {
		name = "sudo"
		args = []string{"-u", n.userName, "bash", "-c", nodeCmd}
	} else {
		name = "bash"
		args = []string{"-c", nodeCmd}
	}

	_, stderr, err = n.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError("nvm", name, args, "", stderr, err)
	}

	return nil
}

func (n *NvmInstaller) configureShell() error {
	nvmDirLine := fmt.Sprintf(`export NVM_DIR="%s/.nvm"`, n.homeDir)
	nvmLoadLine := `[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm`
	nvmCompletionLine := `[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion`
	
	configs := []string{".bashrc", ".zshrc"}
	for _, config := range configs {
		configPath := filepath.Join(n.homeDir, config)
		
		// Skip if file doesn't exist
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		// Read content
		content, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", config, err)
		}

		if strings.Contains(string(content), "nvm.sh") {
			continue
		}

		// Append lines
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open %s for appending: %w", config, err)
		}

		lines := fmt.Sprintf("\n%s\n%s\n%s\n", nvmDirLine, nvmLoadLine, nvmCompletionLine)
		if _, err := f.WriteString(lines); err != nil {
			f.Close()
			return fmt.Errorf("failed to write to %s: %w", config, err)
		}
		f.Close()
	}
	return nil
}
