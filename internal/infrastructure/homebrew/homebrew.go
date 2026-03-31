package homebrew

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/so-install/internal/core/domain"
)

const (
	homebrewScriptURL = "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh"
	brewPath          = "/home/linuxbrew/.linuxbrew/bin/brew"
)

// HomebrewInstaller installs Homebrew and configures the environment.
type HomebrewInstaller struct {
	executor domain.Executor
	homeDir  string
	brewPath string
	userName string
}

// NewHomebrewInstaller creates a new HomebrewInstaller.
func NewHomebrewInstaller(executor domain.Executor) *HomebrewInstaller {
	return &HomebrewInstaller{
		executor: executor,
		homeDir:  domain.GetActualHome(),
		brewPath: brewPath,
		userName: domain.GetActualUser(),
	}
}

var _ domain.SoftwareInstaller = (*HomebrewInstaller)(nil)

// ID returns the SoftwareID for Homebrew.
func (h *HomebrewInstaller) ID() domain.SoftwareID { return domain.Homebrew }

// IsInstalled checks if Homebrew binary exists in the specified location.
func (h *HomebrewInstaller) IsInstalled() (bool, error) {
	_, err := os.Stat(h.brewPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Install installs Homebrew and its dependencies.
func (h *HomebrewInstaller) Install() error {
	// 1. Install system dependencies (must be root)
	deps := []string{"build-essential", "procps", "curl", "file", "git"}
	installDepsCmd := append([]string{"install", "-y"}, deps...)
	_, stderr, err := h.executor.Execute("apt", installDepsCmd...)
	if err != nil {
		return domain.WrapInstallError("homebrew", "apt", installDepsCmd, "", stderr, err)
	}

	// 2. Run official installation script (non-interactive) AS THE ACTUAL USER
	// Homebrew installation script explicitly forbids running as root.
	installCmd := fmt.Sprintf("/bin/bash -c \"$(curl -fsSL %s)\" \"\" --noninteractive", homebrewScriptURL)
	
	var name string
	var args []string
	if h.userName != "" && h.userName != "root" {
		name = "sudo"
		args = []string{"-u", h.userName, "bash", "-c", installCmd}
	} else {
		name = "bash"
		args = []string{"-c", installCmd}
	}

	_, stderr, err = h.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError("homebrew", name, args, "", stderr, err)
	}

	// 3. Configure shell environment
	return h.configureShell()
}

func (h *HomebrewInstaller) configureShell() error {
	shellConfigLine := `eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"`
	configs := []string{".bashrc", ".zshrc"}

	for _, config := range configs {
		configPath := filepath.Join(h.homeDir, config)
		
		// Check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		// Read content to check if line already exists
		content, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", config, err)
		}

		if strings.Contains(string(content), shellConfigLine) {
			continue
		}

		// Append line
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open %s for appending: %w", config, err)
		}

		if _, err := f.WriteString("\n" + shellConfigLine + "\n"); err != nil {
			f.Close()
			return fmt.Errorf("failed to write to %s: %w", config, err)
		}
		f.Close()
	}

	return nil
}
