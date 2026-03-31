package mocks_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMockSoftwareInstaller_ID(t *testing.T) {
	mock := &mocks.MockSoftwareInstaller{
		SoftwareID: domain.SoftwareID("test-software"),
	}
	assert.Equal(t, domain.SoftwareID("test-software"), mock.ID())
}
