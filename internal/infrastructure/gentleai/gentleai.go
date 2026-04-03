package gentleai

import (
	"github.com/so-install/internal/core/domain"
)

const (
	gentleaiInstallCmd = "curl -fsSL https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install.sh | bash"
)

// GentleAIInstaller installs Gentle-AI using the official install script.
type GentleAIInstaller struct {
	executor domain.Executor
	userName string
}

// NewGentleAIInstaller creates a new GentleAIInstaller.
func NewGentleAIInstaller(executor domain.Executor) *GentleAIInstaller {
	return &GentleAIInstaller{
		executor: executor,
		userName: domain.GetActualUser(),
	}
}

var _ domain.SoftwareInstaller = (*GentleAIInstaller)(nil)

// ID returns the SoftwareID for Gentle-AI.
func (g *GentleAIInstaller) ID() domain.SoftwareID { return domain.GentleAI }

// IsInstalled checks if gentle-ai is already installed.
func (g *GentleAIInstaller) IsInstalled() (bool, error) {
	name, args := g.getCommand("gentle-ai", "--version")
	_, _, err := g.executor.Execute(name, args...)
	return err == nil, nil
}

// Install installs Gentle-AI using the official install script.
func (g *GentleAIInstaller) Install() error {
	name, args := g.getCommand("bash", "-c", gentleaiInstallCmd)
	_, stderr, err := g.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError("gentle-ai", name, args, "", stderr, err)
	}
	return nil
}

// getCommand returns the command and arguments, wrapping with sudo -u if userName is provided.
func (g *GentleAIInstaller) getCommand(args ...string) (string, []string) {
	if g.userName != "" && g.userName != "root" {
		return "sudo", append([]string{"-u", g.userName}, args...)
	}
	return args[0], args[1:]
}
