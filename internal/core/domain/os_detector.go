package domain

// OSInfo holds information about the detected operating system.
type OSInfo struct {
	ID        string // e.g. "debian"
	VersionID string // e.g. "12"
}

// OSDetector abstracts OS detection for testability.
type OSDetector interface {
	Detect() (*OSInfo, error)
}
