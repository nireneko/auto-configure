package domain

import (
	"os"
	"os/user"
)

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
