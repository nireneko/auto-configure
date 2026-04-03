package ollama

import (
	"github.com/so-install/internal/core/domain"
)

const (
	ollamaInstallCmd = "curl -fsSL https://ollama.com/install.sh | sh"
)

// OllamaInstaller installs Ollama using the official install script.
type OllamaInstaller struct {
	executor domain.Executor
}

// NewOllamaInstaller creates a new OllamaInstaller.
func NewOllamaInstaller(executor domain.Executor) *OllamaInstaller {
	return &OllamaInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*OllamaInstaller)(nil)

// ID returns the SoftwareID for Ollama.
func (o *OllamaInstaller) ID() domain.SoftwareID { return domain.Ollama }

// IsInstalled checks if ollama is already installed.
func (o *OllamaInstaller) IsInstalled() (bool, error) {
	_, _, err := o.executor.Execute("ollama", "--version")
	return err == nil, nil
}

// Install installs Ollama using the official install script (runs as root — no sudo -u).
func (o *OllamaInstaller) Install() error {
	_, stderr, err := o.executor.Execute("bash", "-c", ollamaInstallCmd)
	if err != nil {
		return domain.WrapInstallError("ollama", "bash", []string{"-c", ollamaInstallCmd}, "", stderr, err)
	}
	return nil
}
