package domain

import (
	"fmt"
	"strings"
)

type OsNotSupportedError struct {
	OS      string
	Version string
}

func (e OsNotSupportedError) Error() string {
	return fmt.Sprintf("OS not supported: %s %s", e.OS, e.Version)
}

type PrivilegeError struct{}

func (e PrivilegeError) Error() string {
	return "insufficient privileges: must be run as root or with sudo"
}

type InstallError struct {
	Software string
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
}

func (e InstallError) Error() string {
	return fmt.Sprintf("failed to install %s: command '%s' exited with code %d", e.Software, e.Command, e.ExitCode)
}

type AptLockError struct {
	InstallError
	LockPath string
}

func (e AptLockError) Error() string {
	return fmt.Sprintf("APT lock held: %s: %s", e.LockPath, e.InstallError.Error())
}

// WrapInstallError wraps a shell error as InstallError or AptLockError.
func WrapInstallError(software, cmd string, args []string, stdout, stderr string, err error) error {
	base := InstallError{
		Software: software,
		Command:  cmd,
		Args:     args,
		Stdout:   stdout,
		Stderr:   stderr,
	}
	if strings.Contains(stderr, "Could not get lock") {
		return AptLockError{InstallError: base}
	}
	return base
}
