package domain_test

import (
	"os"
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

func TestGetActualUID(t *testing.T) {
	t.Run("SUDO_UID set to valid integer returns that value", func(t *testing.T) {
		t.Setenv("SUDO_UID", "1000")
		assert.Equal(t, 1000, domain.GetActualUID())
	})

	t.Run("SUDO_UID not set returns os.Getuid()", func(t *testing.T) {
		t.Setenv("SUDO_UID", "")
		assert.Equal(t, os.Getuid(), domain.GetActualUID())
	})

	t.Run("SUDO_UID set to non-integer returns os.Getuid()", func(t *testing.T) {
		t.Setenv("SUDO_UID", "abc")
		assert.Equal(t, os.Getuid(), domain.GetActualUID())
	})
}

func TestGetActualGID(t *testing.T) {
	t.Run("SUDO_GID set to valid integer returns that value", func(t *testing.T) {
		t.Setenv("SUDO_GID", "1000")
		assert.Equal(t, 1000, domain.GetActualGID())
	})

	t.Run("SUDO_GID not set returns os.Getgid()", func(t *testing.T) {
		t.Setenv("SUDO_GID", "")
		assert.Equal(t, os.Getgid(), domain.GetActualGID())
	})

	t.Run("SUDO_GID set to non-integer returns os.Getgid()", func(t *testing.T) {
		t.Setenv("SUDO_GID", "abc")
		assert.Equal(t, os.Getgid(), domain.GetActualGID())
	})
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
