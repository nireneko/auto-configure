package antigravity_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/antigravity"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAntigravityInstaller_ID(t *testing.T) {
	inst := antigravity.NewAntigravityInstaller(nil)
	assert.Equal(t, domain.Antigravity, inst.ID())
}

func TestAntigravityInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("/usr/bin/agy", "", nil)
	inst := antigravity.NewAntigravityInstaller(m)
	
	installed, err := inst.IsInstalled()
	
	assert.NoError(t, err)
	assert.True(t, installed)
	assert.Equal(t, 1, len(m.Calls))
	assert.Equal(t, "which", m.Calls[0].Name)
	assert.Equal(t, []string{"agy"}, m.Calls[0].Args)
}

func TestAntigravityInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("exit 1"))
	inst := antigravity.NewAntigravityInstaller(m)
	
	installed, err := inst.IsInstalled()
	
	assert.NoError(t, err)
	assert.False(t, installed)
}

func TestAntigravityInstaller_Install_Success(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil) // mkdir
	m.AddResponse("", "", nil) // curl
	m.AddResponse("", "", nil) // sh -c repo entry
	m.AddResponse("", "", nil) // apt update
	m.AddResponse("", "", nil) // apt install
	inst := antigravity.NewAntigravityInstaller(m)
	
	err := inst.Install()
	
	assert.NoError(t, err)
	assert.Equal(t, 5, len(m.Calls))
	
	// mkdir
	assert.Equal(t, "mkdir", m.Calls[0].Name)
	assert.Equal(t, []string{"-p", "/etc/apt/keyrings"}, m.Calls[0].Args)
	
	// curl
	assert.Equal(t, "curl", m.Calls[1].Name)
	assert.Contains(t, m.Calls[1].Args, "https://us-central1-apt.pkg.dev/doc/repo-signing-key.gpg")
	
	// sh -c
	assert.Equal(t, "sh", m.Calls[2].Name)
	assert.Equal(t, "-c", m.Calls[2].Args[0])
	assert.Contains(t, m.Calls[2].Args[1], "antigravity.list")
	
	// apt update
	assert.Equal(t, "apt", m.Calls[3].Name)
	assert.Equal(t, []string{"update"}, m.Calls[3].Args)
	
	// apt install
	assert.Equal(t, "apt", m.Calls[4].Name)
	assert.Equal(t, []string{"install", "-y", "antigravity"}, m.Calls[4].Args)
}

func TestAntigravityInstaller_Install_Failure(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "permission denied", errors.New("exit 1"))
	inst := antigravity.NewAntigravityInstaller(m)
	
	err := inst.Install()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "antigravity")
}
