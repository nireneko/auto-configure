package ddev_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/ddev"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDdev_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := ddev.NewDdevInstaller(exec)
	assert.Equal(t, domain.Ddev, installer.ID())
}
