package domain

import "fmt"

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
	Browser  string
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
}

func (e InstallError) Error() string {
	return fmt.Sprintf("failed to install %s: command '%s' exited with code %d", e.Browser, e.Command, e.ExitCode)
}

type AptLockError struct {
	InstallError
	LockPath string
}

func (e AptLockError) Error() string {
	return fmt.Sprintf("APT lock held: %s: %s", e.LockPath, e.InstallError.Error())
}
