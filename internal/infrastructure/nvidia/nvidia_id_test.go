package nvidia_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/nvidia"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNvidiaInstaller_ID(t *testing.T) {
	inst := nvidia.NewNvidiaInstaller(&mocks.MockExecutor{}, &domain.OSInfo{})
	assert.Equal(t, domain.NvidiaDrivers, inst.ID())
}
