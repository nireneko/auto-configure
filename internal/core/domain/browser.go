package domain

// BrowserID identifies a supported browser.
type BrowserID string

const (
	Brave    BrowserID = "brave"
	Firefox  BrowserID = "firefox"
	Chrome   BrowserID = "chrome"
	Chromium BrowserID = "chromium"
)

// AllBrowsers returns all supported browsers in display order.
func AllBrowsers() []BrowserID {
	return []BrowserID{Brave, Firefox, Chrome, Chromium}
}

// DisplayName returns a human-readable name for the browser.
func (b BrowserID) DisplayName() string {
	switch b {
	case Brave:
		return "Brave"
	case Firefox:
		return "Firefox"
	case Chrome:
		return "Google Chrome"
	case Chromium:
		return "Chromium"
	default:
		return string(b)
	}
}

// BrowserInstaller handles installation of a specific browser.
type BrowserInstaller interface {
	Install() error
	IsInstalled() (bool, error)
	ID() BrowserID
}

// InstallResult holds the outcome of a browser installation attempt.
type InstallResult struct {
	Browser          BrowserID
	AlreadyInstalled bool
	Err              error
}
