package nvm_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/nvm"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNvm_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := nvm.NewNvmInstaller(exec)
	assert.Equal(t, domain.Nvm, installer.ID())
}
