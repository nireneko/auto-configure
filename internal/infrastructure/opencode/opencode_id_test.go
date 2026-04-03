package opencode_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/opencode"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOpenCode_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := opencode.NewOpenCodeInstaller(exec)
	assert.Equal(t, domain.OpenCode, installer.ID())
}
