package usecases

import (
	"strings"
	"time"

	"github.com/so-install/internal/core/domain"
)

const maxAptLockRetries = 3

// InstallBrowsersUseCase orchestrates browser installations sequentially.
type InstallBrowsersUseCase struct {
	installers map[domain.BrowserID]domain.BrowserInstaller
	sleepFn    func(time.Duration)
}

// NewInstallBrowsersUseCase creates a new use case.
// In production, pass time.Sleep as sleepFn.
func NewInstallBrowsersUseCase(
	installers map[domain.BrowserID]domain.BrowserInstaller,
	sleepFn func(time.Duration),
) *InstallBrowsersUseCase {
	return &InstallBrowsersUseCase{installers: installers, sleepFn: sleepFn}
}

// Execute installs the selected browsers sequentially and returns results.
func (uc *InstallBrowsersUseCase) Execute(selected []domain.BrowserID) []domain.InstallResult {
	results := make([]domain.InstallResult, 0, len(selected))
	for _, id := range selected {
		installer := uc.installers[id]
		result := uc.installOne(id, installer)
		results = append(results, result)
	}
	return results
}

func (uc *InstallBrowsersUseCase) installOne(id domain.BrowserID, installer domain.BrowserInstaller) domain.InstallResult {
	installed, err := installer.IsInstalled()
	if err == nil && installed {
		return domain.InstallResult{Browser: id, AlreadyInstalled: true}
	}

	var lastErr error
	for attempt := 0; attempt < maxAptLockRetries; attempt++ {
		lastErr = installer.Install()
		if lastErr == nil {
			return domain.InstallResult{Browser: id}
		}
		if !isAptLockError(lastErr) {
			return domain.InstallResult{Browser: id, Err: lastErr}
		}
		if attempt < maxAptLockRetries-1 {
			uc.sleepFn(5 * time.Second)
		}
	}
	return domain.InstallResult{Browser: id, Err: lastErr}
}

func isAptLockError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(domain.AptLockError); ok {
		return true
	}
	return strings.Contains(err.Error(), "Could not get lock")
}
