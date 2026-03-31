package docker_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/docker"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDocker_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	installer := docker.NewDockerInstaller(exec, "testuser")
	assert.Equal(t, domain.Docker, installer.ID())
}
