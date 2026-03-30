package flatpak

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFlatpakAppInstaller_Install(t *testing.T) {
	t.Run("Successful installation", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		installer := NewFlatpakAppInstaller(executor, "com.bitwarden.desktop", domain.Bitwarden)

		err := installer.Install()

		assert.NoError(t, err)
		assert.Equal(t, "flatpak", executor.Calls[0].Name)
		assert.Equal(t, []string{"install", "-y", "flathub", "com.bitwarden.desktop"}, executor.Calls[0].Args)
	})

	t.Run("Failed installation", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "error installing", errors.New("exit status 1"))
		installer := NewFlatpakAppInstaller(executor, "com.bitwarden.desktop", domain.Bitwarden)

		err := installer.Install()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to install bitwarden")
		assert.Contains(t, err.Error(), "error installing")
	})
}

func TestFlatpakAppInstaller_IsInstalled(t *testing.T) {
	t.Run("App is installed", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("Bitwarden - The open source password manager", "", nil)
		installer := NewFlatpakAppInstaller(executor, "com.bitwarden.desktop", domain.Bitwarden)

		installed, err := installer.IsInstalled()

		assert.NoError(t, err)
		assert.True(t, installed)
		assert.Equal(t, "flatpak", executor.Calls[0].Name)
		assert.Equal(t, []string{"info", "com.bitwarden.desktop"}, executor.Calls[0].Args)
	})

	t.Run("App is not installed", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "error: app not found", errors.New("exit status 1"))
		installer := NewFlatpakAppInstaller(executor, "com.bitwarden.desktop", domain.Bitwarden)

		installed, err := installer.IsInstalled()

		assert.NoError(t, err)
		assert.False(t, installed)
	})
}

func TestFlatpakAppInstaller_ID(t *testing.T) {
	installer := NewFlatpakAppInstaller(nil, "com.bitwarden.desktop", domain.Bitwarden)
	assert.Equal(t, domain.Bitwarden, installer.ID())
}
