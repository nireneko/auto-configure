package npm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/so-install/internal/core/domain"
)

// NpmInstaller installs software via npm global install.
type NpmInstaller struct {
	executor    domain.Executor
	packageName string
	binaryName  string
	softwareID  domain.SoftwareID
	homeDir     string
}

// NewNpmInstaller creates a new NpmInstaller.
func NewNpmInstaller(executor domain.Executor, packageName, binaryName string, id domain.SoftwareID) *NpmInstaller {
	home, _ := os.UserHomeDir()
	return &NpmInstaller{
		executor:    executor,
		packageName: packageName,
		binaryName:  binaryName,
		softwareID:  id,
		homeDir:     home,
	}
}

var _ domain.SoftwareInstaller = (*NpmInstaller)(nil)

// ID returns the SoftwareID.
func (n *NpmInstaller) ID() domain.SoftwareID { return n.softwareID }

// IsInstalled checks if the binary exists by running it with --version.
func (n *NpmInstaller) IsInstalled() (bool, error) {
	name, args := n.getCommand(n.binaryName, "--version")
	_, _, err := n.executor.Execute(name, args...)
	return err == nil, nil
}

// Install installs the npm package globally.
func (n *NpmInstaller) Install() error {
	name, args := n.getCommand("npm", "install", "-g", n.packageName)
	_, stderr, err := n.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError(string(n.softwareID), name, args, "", stderr, err)
	}
	return nil
}

// getCommand returns the command and arguments, prepending nvm source if nvm is found.
func (n *NpmInstaller) getCommand(args ...string) (string, []string) {
	nvmScript := filepath.Join(n.homeDir, ".nvm", "nvm.sh")
	if _, err := os.Stat(nvmScript); err == nil {
		// Prepend source nvm.sh and join with &&
		fullCmd := fmt.Sprintf("source %s && %s", nvmScript, strings.Join(args, " "))
		return "bash", []string{"-c", fullCmd}
	}
	// Fallback to direct command
	return args[0], args[1:]
}
