package domain_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestErrors_ErrorStrings(t *testing.T) {
	// OsNotSupportedError
	osErr := domain.OsNotSupportedError{OS: "Windows", Version: "10"}
	assert.Equal(t, "OS not supported: Windows 10", osErr.Error())

	// PrivilegeError
	privErr := domain.PrivilegeError{}
	assert.Equal(t, "insufficient privileges: must be run as root or with sudo", privErr.Error())

	// InstallError without stderr
	instErrNoStderr := domain.InstallError{
		Software: "test-soft",
		Command:  "apt-get",
	}
	assert.Equal(t, "failed to install test-soft: command 'apt-get' failed", instErrNoStderr.Error())

	// InstallError with stderr
	instErrWithStderr := domain.InstallError{
		Software: "test-soft",
		Command:  "apt-get",
		Stderr:   "not found",
	}
	assert.Equal(t, "failed to install test-soft: command 'apt-get' failed: not found", instErrWithStderr.Error())

	// AptLockError
	aptErr := domain.AptLockError{
		InstallError: instErrWithStderr,
		LockPath:     "/var/lib/dpkg/lock",
	}
	assert.Equal(t, "APT lock held: /var/lib/dpkg/lock: failed to install test-soft: command 'apt-get' failed: not found", aptErr.Error())
}

func TestWrapInstallError(t *testing.T) {
	err := errors.New("some error")

	// Standard error wrap
	wrapped := domain.WrapInstallError("soft", "cmd", []string{"arg"}, "out", "err out", err)
	assert.IsType(t, domain.InstallError{}, wrapped)

	// Apt lock error wrap
	wrappedApt := domain.WrapInstallError("soft", "cmd", nil, "", "Could not get lock /var/lib", err)
	assert.IsType(t, domain.AptLockError{}, wrappedApt)
}
