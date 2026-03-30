package domain

// SoftwareID identifies installable software.
type SoftwareID string

const (
	SystemUpdate SoftwareID = "system-update"
	BaseDeps     SoftwareID = "base-deps"
	Brave    SoftwareID = "brave"
	Firefox  SoftwareID = "firefox"
	Chrome   SoftwareID = "chrome"
	Chromium SoftwareID = "chromium"
	Docker   SoftwareID = "docker"
	Ddev     SoftwareID = "ddev"
	OpenVpn  SoftwareID = "openvpn"
	Nvm      SoftwareID = "nvm"
	Gemini   SoftwareID = "gemini"
	ClaudeCode SoftwareID = "claude"
	Flatpak  SoftwareID = "flatpak"
	Bitwarden SoftwareID = "bitwarden"
	Homebrew  SoftwareID = "homebrew"
)

// InstallStep defines a group of software to be installed together.
type InstallStep struct {
	ID       string
	Software []SoftwareID
	Critical bool
}

// GetSteps returns the predefined installation phases.
func GetSteps() []InstallStep {
	return []InstallStep{
		{
			ID:       "system-prep",
			Software: []SoftwareID{SystemUpdate, BaseDeps},
			Critical: true,
		},
		{
			ID:       "browsers",
			Software: []SoftwareID{Brave, Firefox, Chrome, Chromium},
			Critical: false,
		},
		{
			ID:       "docker",
			Software: []SoftwareID{Docker},
			Critical: true,
		},
		{
			ID:       "ddev",
			Software: []SoftwareID{Ddev},
			Critical: false,
		},
		{
			ID:       "openvpn",
			Software: []SoftwareID{OpenVpn},
			Critical: false,
		},
		{
			ID:       "nvm",
			Software: []SoftwareID{Nvm},
			Critical: false,
		},
		{
			ID:       "ai-cli",
			Software: []SoftwareID{Gemini, ClaudeCode},
			Critical: false,
		},
		{
			ID:       "flatpak",
			Software: []SoftwareID{Flatpak},
			Critical: false,
		},
		{
			ID:       "apps",
			Software: []SoftwareID{Bitwarden, Homebrew},
			Critical: false,
		},
	}
}

// AllSoftware returns all supported software in display order.
func AllSoftware() []SoftwareID {
	return []SoftwareID{Brave, Firefox, Chrome, Chromium, Docker, Ddev, OpenVpn, Nvm, Gemini, ClaudeCode, Flatpak, Bitwarden, Homebrew}
}

// DisplayName returns a human-readable name for the software.
func (s SoftwareID) DisplayName() string {
	switch s {
	case SystemUpdate:
		return "System Update"
	case BaseDeps:
		return "Base Dependencies"
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
	case Ddev:
		return "DDEV"
	case OpenVpn:
		return "OpenVPN"
	case Nvm:
		return "NVM & NPM"
	case Gemini:
		return "Google Gemini CLI"
	case ClaudeCode:
		return "Claude Code (Anthropic)"
	case Flatpak:
		return "Flatpak"
	case Bitwarden:
		return "Bitwarden"
	case Homebrew:
		return "Homebrew"
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
