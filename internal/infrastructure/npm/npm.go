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
	userName    string
}

// NewNpmInstaller creates a new NpmInstaller.
func NewNpmInstaller(executor domain.Executor, packageName, binaryName string, id domain.SoftwareID) *NpmInstaller {
	return &NpmInstaller{
		executor:    executor,
		packageName: packageName,
		binaryName:  binaryName,
		softwareID:  id,
		homeDir:     domain.GetActualHome(),
		userName:    domain.GetActualUser(),
	}
}

var _ domain.SoftwareInstaller = (*NpmInstaller)(nil)

// ID returns the SoftwareID.
func (n *NpmInstaller) ID() domain.SoftwareID { return n.softwareID }

// IsInstalled checks if the binary exists by running it with --version.
func (n *NpmInstaller) IsInstalled() (bool, error) {
	name, args := n.getCommand(n.userName, n.binaryName, "--version")
	_, _, err := n.executor.Execute(name, args...)
	return err == nil, nil
}

// Install installs the npm package globally.
func (n *NpmInstaller) Install() error {
	name, args := n.getCommand(n.userName, "npm", "install", "-g", n.packageName)
	_, stderr, err := n.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError(string(n.softwareID), name, args, "", stderr, err)
	}
	return nil
}

// getCommand returns the command and arguments, prepending nvm source if nvm is found.
// It also wraps the command with sudo -u if userName is provided and not root.
func (n *NpmInstaller) getCommand(userName string, args ...string) (string, []string) {
	nvmScript := filepath.Join(n.homeDir, ".nvm", "nvm.sh")
	
	var finalCmd string
	if _, err := os.Stat(nvmScript); err == nil {
		// Prepend source nvm.sh and join with &&
		finalCmd = fmt.Sprintf("source %s && %s", nvmScript, strings.Join(args, " "))
	} else {
		finalCmd = strings.Join(args, " ")
	}

	if userName != "" && userName != "root" {
		return "sudo", []string{"-u", userName, "bash", "-c", finalCmd}
	}
	
	return "bash", []string{"-c", finalCmd}
}
