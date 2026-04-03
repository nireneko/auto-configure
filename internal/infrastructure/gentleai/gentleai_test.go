package gentleai

import (
	"fmt"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewGentleAIInstaller(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := NewGentleAIInstaller(executor)
	assert.NotNil(t, installer)
	assert.Equal(t, executor, installer.executor)
}

func TestGentleAIInstaller_ID(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &GentleAIInstaller{executor: executor, userName: "testuser"}
	assert.Equal(t, domain.GentleAI, installer.ID())
}

func TestGentleAIInstaller_IsInstalled_True(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("gentle-ai version 0.1.0", "", nil)
	installer := &GentleAIInstaller{executor: executor, userName: ""}

	installed, err := installer.IsInstalled()

	assert.NoError(t, err)
	assert.True(t, installed)
	assert.Equal(t, 1, len(executor.Calls))
	assert.Equal(t, "gentle-ai", executor.Calls[0].Name)
	assert.Equal(t, "--version", executor.Calls[0].Args[0])
}

func TestGentleAIInstaller_IsInstalled_False(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "not found", fmt.Errorf("exit status 127"))
	installer := &GentleAIInstaller{executor: executor, userName: ""}

	installed, err := installer.IsInstalled()

	assert.NoError(t, err)
	assert.False(t, installed)
}

func TestGentleAIInstaller_Install_SudoUser(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &GentleAIInstaller{executor: executor, userName: "testuser"}

	err := installer.Install()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(executor.Calls))
	call := executor.Calls[0]
	assert.Equal(t, "sudo", call.Name)
	assert.Equal(t, "-u", call.Args[0])
	assert.Equal(t, "testuser", call.Args[1])
	assert.Equal(t, "bash", call.Args[2])
	assert.Equal(t, "-c", call.Args[3])
	assert.Contains(t, call.Args[4], "curl -fsSL")
	assert.Contains(t, call.Args[4], "gentle-ai/main/scripts/install.sh")
}

func TestGentleAIInstaller_Install_RootUser(t *testing.T) {
	executor := &mocks.MockExecutor{}
	installer := &GentleAIInstaller{executor: executor, userName: "root"}

	err := installer.Install()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(executor.Calls))
	call := executor.Calls[0]
	assert.Equal(t, "bash", call.Name)
	assert.Equal(t, "-c", call.Args[0])
	assert.Contains(t, call.Args[1], "curl -fsSL")
}

func TestGentleAIInstaller_Install_Error(t *testing.T) {
	executor := &mocks.MockExecutor{}
	executor.AddResponse("", "curl: network error", fmt.Errorf("exit status 1"))
	installer := &GentleAIInstaller{executor: executor, userName: ""}

	err := installer.Install()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to install gentle-ai")
}
