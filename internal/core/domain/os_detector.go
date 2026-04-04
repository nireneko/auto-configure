package domain

// DesktopEnvironment identifies the user's desktop environment.
type DesktopEnvironment string

const (
	KDE   DesktopEnvironment = "kde"
	GNOME DesktopEnvironment = "gnome"
	Other DesktopEnvironment = "other"
)

// OSInfo holds information about the detected operating system.
type OSInfo struct {
	ID                 string // e.g. "debian"
	VersionID          string // e.g. "12"
	DesktopEnvironment DesktopEnvironment
	IsWayland          bool // true when XDG_SESSION_TYPE==wayland or WAYLAND_DISPLAY is set
}

// OSDetector abstracts OS detection for testability.
type OSDetector interface {
	Detect() (*OSInfo, error)
}
