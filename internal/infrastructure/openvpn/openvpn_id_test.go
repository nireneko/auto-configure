package openvpn_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/openvpn"
	"github.com/so-install/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOpenVPN_ID(t *testing.T) {
	exec := &mocks.MockExecutor{}
	osInfo := &domain.OSInfo{ID: "debian"}
	installer := openvpn.NewOpenVpnInstaller(exec, osInfo)
	assert.Equal(t, domain.OpenVpn, installer.ID())
}
