package homebrew_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/homebrew"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHomebrew_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := homebrew.NewHomebrewInstaller(exec)
	assert.Equal(t, domain.Homebrew, installer.ID())
}
