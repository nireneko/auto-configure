package npm_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/npm"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNpm_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := npm.NewNpmInstaller(exec, "test-pkg", "test-bin", domain.SystemUpdate)
	assert.Equal(t, domain.SystemUpdate, installer.ID())
}
