package domain_test

import (
	"os/user"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetActualUser(t *testing.T) {
	// Case 1: SUDO_USER is set
	t.Setenv("SUDO_USER", "test_sudo_user")
	assert.Equal(t, "test_sudo_user", domain.GetActualUser())

	// Case 2: SUDO_USER is not set
	t.Setenv("SUDO_USER", "")
	u, err := user.Current()
	if err == nil && u != nil {
		assert.Equal(t, u.Username, domain.GetActualUser())
	}
}

func TestGetActualHome(t *testing.T) {
	// Case 1: SUDO_USER is set to a user that exists (root should exist on linux)
	// We use "root" to ensure user.Lookup doesn't fail on standard systems.
	t.Setenv("SUDO_USER", "root")
	home := domain.GetActualHome()
	assert.NotEmpty(t, home)

	// Case 2: SUDO_USER is not set
	t.Setenv("SUDO_USER", "")
	u, err := user.Current()
	if err == nil && u != nil {
		assert.Equal(t, u.HomeDir, domain.GetActualHome())
	}
}
