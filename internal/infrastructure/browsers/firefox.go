package browsers

import (
	"github.com/so-install/internal/core/domain"
)

const (
	firefoxGPGURL      = "https://packages.mozilla.org/apt/repo-signing-key.gpg"
	firefoxGPGPath     = "/etc/apt/keyrings/packages.mozilla.org.gpg"
	firefoxSourceDeb   = "deb [signed-by=/etc/apt/keyrings/packages.mozilla.org.gpg] https://packages.mozilla.org/apt mozilla main"
	firefoxSourcePath  = "/etc/apt/sources.list.d/mozilla.list"
	firefoxPrefPath    = "/etc/apt/preferences.d/mozilla"
	firefoxPrefContent = "Package: *\nPin: origin packages.mozilla.org\nPin-Priority: 1000\n"
)

// FirefoxInstaller installs Firefox from the official Mozilla repository.
type FirefoxInstaller struct {
	executor domain.Executor
}

// NewFirefoxInstaller creates a new FirefoxInstaller.
func NewFirefoxInstaller(executor domain.Executor) *FirefoxInstaller {
	return &FirefoxInstaller{executor: executor}
}

var _ domain.BrowserInstaller = (*FirefoxInstaller)(nil)

// ID returns the BrowserID for Firefox.
func (f *FirefoxInstaller) ID() domain.BrowserID { return domain.Firefox }

// IsInstalled checks if firefox is already installed.
func (f *FirefoxInstaller) IsInstalled() (bool, error) {
	_, _, err := f.executor.Execute("which", "firefox")
	return err == nil, nil
}

// Install installs Firefox from the Mozilla official APT repository.
func (f *FirefoxInstaller) Install() error {
	steps := [][]string{
		{"wget", "-q", firefoxGPGURL, "-O", firefoxGPGPath},
		{"sh", "-c", "echo '" + firefoxSourceDeb + "' | tee " + firefoxSourcePath},
		{"sh", "-c", "printf '" + firefoxPrefContent + "' | tee " + firefoxPrefPath},
		{"apt", "update"},
		{"apt", "install", "-y", "firefox"},
	}
	for _, step := range steps {
		_, stderr, err := f.executor.Execute(step[0], step[1:]...)
		if err != nil {
			return wrapInstallError("firefox", step[0], step[1:], "", stderr, err)
		}
	}
	return nil
}
