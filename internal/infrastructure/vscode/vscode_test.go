package vscode_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/vscode"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestVsCodeInstaller_ID(t *testing.T) {
	inst := vscode.NewVsCodeInstaller(nil)
	assert.Equal(t, domain.VsCode, inst.ID())
}

func TestVsCodeInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("/usr/bin/code", "", nil)
	inst := vscode.NewVsCodeInstaller(m)
	
	installed, err := inst.IsInstalled()
	
	assert.NoError(t, err)
	assert.True(t, installed)
	assert.Equal(t, 1, len(m.Calls))
	assert.Equal(t, "which", m.Calls[0].Name)
	assert.Equal(t, []string{"code"}, m.Calls[0].Args)
}

func TestVsCodeInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("exit 1"))
	inst := vscode.NewVsCodeInstaller(m)
	
	installed, err := inst.IsInstalled()
	
	assert.NoError(t, err)
	assert.False(t, installed)
}

func TestVsCodeInstaller_Install_Success(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil) // wget
	m.AddResponse("", "", nil) // apt install
	inst := vscode.NewVsCodeInstaller(m)
	
	err := inst.Install()
	
	assert.NoError(t, err)
	assert.Equal(t, 2, len(m.Calls))
	
	// wget call
	assert.Equal(t, "wget", m.Calls[0].Name)
	assert.Contains(t, m.Calls[0].Args, "https://go.microsoft.com/fwlink/?LinkID=760868")
	assert.Contains(t, m.Calls[0].Args, "/tmp/vscode.deb")
	
	// apt install call
	assert.Equal(t, "apt", m.Calls[1].Name)
	assert.Equal(t, []string{"install", "-y", "/tmp/vscode.deb"}, m.Calls[1].Args)
}

func TestVsCodeInstaller_Install_WgetFailure(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "connection refused", errors.New("exit 1"))
	inst := vscode.NewVsCodeInstaller(m)
	
	err := inst.Install()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vscode") // WrapInstallError wraps it with "vscode"
}
