package gitlab_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/gitlab"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGitlabConfigurator_IDAndIsInstalled(t *testing.T) {
	exec := &mocks.MockExecutor{}
	conf := gitlab.NewGitlabTokenConfigurator(exec)

	assert.Equal(t, domain.GitlabTokenConfig, conf.ID())

	installed, err := conf.IsInstalled()
	assert.NoError(t, err)
	assert.False(t, installed)
}

func TestGitlabConfigurator_EmptyToken(t *testing.T) {
	exec := &mocks.MockExecutor{}
	conf := gitlab.NewGitlabTokenConfigurator(exec)
	// No token set
	err := conf.Install()
	assert.ErrorContains(t, err, "gitlab token is not set")
}
