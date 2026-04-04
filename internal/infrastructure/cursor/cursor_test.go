package cursor_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/cursor"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCursorInstaller_ID(t *testing.T) {
	inst := cursor.NewCursorInstaller(nil)
	assert.Equal(t, domain.Cursor, inst.ID())
}

func TestCursorInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("/usr/bin/cursor", "", nil)
	inst := cursor.NewCursorInstaller(m)
	
	installed, err := inst.IsInstalled()
	
	assert.NoError(t, err)
	assert.True(t, installed)
	assert.Equal(t, 1, len(m.Calls))
	assert.Equal(t, "which", m.Calls[0].Name)
	assert.Equal(t, []string{"cursor"}, m.Calls[0].Args)
}

func TestCursorInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("exit 1"))
	inst := cursor.NewCursorInstaller(m)
	
	installed, err := inst.IsInstalled()
	
	assert.NoError(t, err)
	assert.False(t, installed)
}

func TestCursorInstaller_Install_Success(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil) // wget
	m.AddResponse("", "", nil) // apt install
	inst := cursor.NewCursorInstaller(m)
	
	err := inst.Install()
	
	assert.NoError(t, err)
	assert.Equal(t, 2, len(m.Calls))
	
	// wget call
	assert.Equal(t, "wget", m.Calls[0].Name)
	assert.Contains(t, m.Calls[0].Args, "https://downloader.cursor.sh/linux/debian/amd64")
	assert.Contains(t, m.Calls[0].Args, "/tmp/cursor.deb")
	
	// apt install call
	assert.Equal(t, "apt", m.Calls[1].Name)
	assert.Equal(t, []string{"install", "-y", "/tmp/cursor.deb"}, m.Calls[1].Args)
}

func TestCursorInstaller_Install_WgetFailure(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "connection refused", errors.New("exit 1"))
	inst := cursor.NewCursorInstaller(m)
	
	err := inst.Install()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cursor")
}
