package flatpak_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/flatpak"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFlatpak_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	osDet := &mocks.MockOSDetector{}
	installer := flatpak.NewFlatpakInstaller(exec, osDet)
	assert.Equal(t, domain.Flatpak, installer.ID())
}
