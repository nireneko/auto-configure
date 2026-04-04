package vscode

import (
	"github.com/so-install/internal/core/domain"
)

const (
	vscodeURL = "https://go.microsoft.com/fwlink/?LinkID=760868"
	vscodeDeb = "/tmp/vscode.deb"
)

// VsCodeInstaller installs Visual Studio Code from the official download.
type VsCodeInstaller struct {
	executor domain.Executor
}

// NewVsCodeInstaller creates a new VsCodeInstaller.
func NewVsCodeInstaller(executor domain.Executor) *VsCodeInstaller {
	return &VsCodeInstaller{executor: executor}
}

var _ domain.SoftwareInstaller = (*VsCodeInstaller)(nil)

// ID returns the SoftwareID for VsCode.
func (v *VsCodeInstaller) ID() domain.SoftwareID { return domain.VsCode }

// IsInstalled checks if code is already installed.
func (v *VsCodeInstaller) IsInstalled() (bool, error) {
	_, _, err := v.executor.Execute("which", "code")
	return err == nil, nil
}

// Install downloads and installs Visual Studio Code.
func (v *VsCodeInstaller) Install() error {
	steps := [][]string{
		{"wget", vscodeURL, "-O", vscodeDeb},
		{"apt", "install", "-y", vscodeDeb},
	}
	for _, step := range steps {
		_, stderr, err := v.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return domain.WrapInstallError("vscode", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
