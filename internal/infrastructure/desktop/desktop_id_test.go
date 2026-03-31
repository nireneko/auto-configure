package desktop_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/desktop"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestScreenLock_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	osDet := &mocks.MockOSDetector{}
	installer := desktop.NewScreenLockInstaller(exec, osDet)
	assert.Equal(t, domain.ScreenLockConfig, installer.ID())
}
