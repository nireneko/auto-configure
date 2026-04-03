package opencode

import (
	"github.com/so-install/internal/core/domain"
)

const (
	openCodeInstallCmd = "curl -fsSL https://opencode.ai/install | bash"
)

// OpenCodeInstaller installs OpenCode using the official install script.
type OpenCodeInstaller struct {
	executor domain.Executor
	userName string
}

// NewOpenCodeInstaller creates a new OpenCodeInstaller using the actual (non-sudo) user.
func NewOpenCodeInstaller(executor domain.Executor) *OpenCodeInstaller {
	return &OpenCodeInstaller{
		executor: executor,
		userName: domain.GetActualUser(),
	}
}

var _ domain.SoftwareInstaller = (*OpenCodeInstaller)(nil)

// ID returns the SoftwareID for OpenCode.
func (o *OpenCodeInstaller) ID() domain.SoftwareID { return domain.OpenCode }

// IsInstalled checks if opencode is already installed.
func (o *OpenCodeInstaller) IsInstalled() (bool, error) {
	_, _, err := o.executor.Execute("opencode", "--version")
	return err == nil, nil
}

// Install installs OpenCode using the official install script.
// Runs as the real user via sudo -u when a non-root user is detected.
func (o *OpenCodeInstaller) Install() error {
	var name string
	var args []string

	if o.userName != "" && o.userName != "root" {
		name = "sudo"
		args = []string{"-u", o.userName, "bash", "-c", openCodeInstallCmd}
	} else {
		name = "bash"
		args = []string{"-c", openCodeInstallCmd}
	}

	_, stderr, err := o.executor.Execute(name, args...)
	if err != nil {
		return domain.WrapInstallError("opencode", name, args, "", stderr, err)
	}
	return nil
}
