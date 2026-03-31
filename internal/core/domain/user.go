package domain

import (
	"os"
	"os/user"
	"strconv"
)

// GetActualUID returns the real user's UID even when running under sudo.
func GetActualUID() int {
	if s := os.Getenv("SUDO_UID"); s != "" {
		if uid, err := strconv.Atoi(s); err == nil {
			return uid
		}
	}
	return os.Getuid()
}

// GetActualGID returns the real user's GID even when running under sudo.
func GetActualGID() int {
	if s := os.Getenv("SUDO_GID"); s != "" {
		if gid, err := strconv.Atoi(s); err == nil {
			return gid
		}
	}
	return os.Getgid()
}

// GetActualUser returns the real user's name even when running under sudo.
func GetActualUser() string {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		return sudoUser
	}
	u, _ := user.Current()
	if u != nil {
		return u.Username
	}
	return ""
}

// GetActualHome returns the real user's home directory even when running under sudo.
func GetActualHome() string {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		if u, err := user.Lookup(sudoUser); err == nil {
			return u.HomeDir
		}
	}
	// Fallback to current user's home
	home, _ := os.UserHomeDir()
	return home
}
