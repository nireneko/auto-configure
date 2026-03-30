package flatpak

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFlatpakInstaller_Install(t *testing.T) {
	tests := []struct {
		name          string
		de            domain.DesktopEnvironment
		expectedCalls []mocks.ExecutorCall
	}{
		{
			name: "KDE installation",
			de:   domain.KDE,
			expectedCalls: []mocks.ExecutorCall{
				{Name: "apt", Args: []string{"update"}},
				{Name: "apt", Args: []string{"install", "-y", "flatpak"}},
				{Name: "flatpak", Args: []string{"remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"}},
				{Name: "apt", Args: []string{"install", "-y", "plasma-discover-backend-flatpak"}},
			},
		},
		{
			name: "GNOME installation",
			de:   domain.GNOME,
			expectedCalls: []mocks.ExecutorCall{
				{Name: "apt", Args: []string{"update"}},
				{Name: "apt", Args: []string{"install", "-y", "flatpak"}},
				{Name: "flatpak", Args: []string{"remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"}},
				{Name: "apt", Args: []string{"install", "-y", "gnome-software-plugin-flatpak"}},
			},
		},
		{
			name: "Other DE installation",
			de:   domain.Other,
			expectedCalls: []mocks.ExecutorCall{
				{Name: "apt", Args: []string{"update"}},
				{Name: "apt", Args: []string{"install", "-y", "flatpak"}},
				{Name: "flatpak", Args: []string{"remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &mocks.MockExecutor{}
			detector := &mocks.MockOSDetector{ReturnID: "debian", ReturnVersionID: "12", ReturnDE: tt.de}
			installer := NewFlatpakInstaller(executor, detector)

			err := installer.Install()

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCalls, executor.Calls)
		})
	}
}

func TestFlatpakInstaller_IsInstalled(t *testing.T) {
	t.Run("Already installed", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("flatpak 1.14.4", "", nil)
		installer := NewFlatpakInstaller(executor, nil)

		installed, err := installer.IsInstalled()

		assert.NoError(t, err)
		assert.True(t, installed)
		assert.Equal(t, "flatpak", executor.Calls[0].Name)
	})

	t.Run("Not installed", func(t *testing.T) {
		executor := &mocks.MockExecutor{}
		executor.AddResponse("", "command not found", errors.New("exit status 127"))
		installer := NewFlatpakInstaller(executor, nil)

		installed, err := installer.IsInstalled()

		assert.NoError(t, err)
		assert.False(t, installed)
	})
}
