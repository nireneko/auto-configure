package domain

// SoftwareID identifies installable software.
type SoftwareID string

const (
	Brave    SoftwareID = "brave"
	Firefox  SoftwareID = "firefox"
	Chrome   SoftwareID = "chrome"
	Chromium SoftwareID = "chromium"
	Docker   SoftwareID = "docker"
)

// AllSoftware returns all supported software in display order.
func AllSoftware() []SoftwareID {
	return []SoftwareID{Brave, Firefox, Chrome, Chromium, Docker}
}

// DisplayName returns a human-readable name for the software.
func (s SoftwareID) DisplayName() string {
	switch s {
	case Brave:
		return "Brave"
	case Firefox:
		return "Firefox"
	case Chrome:
		return "Google Chrome"
	case Chromium:
		return "Chromium"
	case Docker:
		return "Docker CE"
	default:
		return string(s)
	}
}

// SoftwareInstaller handles installation of specific software.
type SoftwareInstaller interface {
	Install() error
	IsInstalled() (bool, error)
	ID() SoftwareID
}

// InstallResult holds the outcome of a software installation attempt.
type InstallResult struct {
	Software         SoftwareID
	AlreadyInstalled bool
	Err              error
}
