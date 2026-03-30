package usecases

import (
	"strings"
	"time"

	"github.com/so-install/internal/core/domain"
)

const maxAptLockRetries = 3

// InstallSoftwareUseCase orchestrates software installations sequentially.
type InstallSoftwareUseCase struct {
	installers map[domain.SoftwareID]domain.SoftwareInstaller
	sleepFn    func(time.Duration)
}

// NewInstallSoftwareUseCase creates a new use case.
// In production, pass time.Sleep as sleepFn.
func NewInstallSoftwareUseCase(
	installers map[domain.SoftwareID]domain.SoftwareInstaller,
	sleepFn func(time.Duration),
) *InstallSoftwareUseCase {
	return &InstallSoftwareUseCase{installers: installers, sleepFn: sleepFn}
}

// Execute installs the selected software sequentially following domain steps.
func (uc *InstallSoftwareUseCase) Execute(selected []domain.SoftwareID) []domain.InstallResult {
	selectedMap := make(map[domain.SoftwareID]bool)
	for _, id := range selected {
		selectedMap[id] = true
	}

	results := make([]domain.InstallResult, 0, len(selected))
	steps := domain.GetSteps()

	for _, step := range steps {
		stepFailed := false
		for _, id := range step.Software {
			if !selectedMap[id] {
				continue
			}

			installer := uc.installers[id]
			result := uc.installOne(id, installer)
			results = append(results, result)

			if result.Err != nil {
				stepFailed = true
			}
		}

		if stepFailed && step.Critical {
			break
		}
	}
	return results
}

func (uc *InstallSoftwareUseCase) installOne(id domain.SoftwareID, installer domain.SoftwareInstaller) domain.InstallResult {
	installed, err := installer.IsInstalled()
	if err == nil && installed {
		return domain.InstallResult{Software: id, AlreadyInstalled: true}
	}

	var lastErr error
	for attempt := 0; attempt < maxAptLockRetries; attempt++ {
		lastErr = installer.Install()
		if lastErr == nil {
			return domain.InstallResult{Software: id}
		}
		if !isAptLockError(lastErr) {
			return domain.InstallResult{Software: id, Err: lastErr}
		}
		if attempt < maxAptLockRetries-1 {
			uc.sleepFn(5 * time.Second)
		}
	}
	return domain.InstallResult{Software: id, Err: lastErr}
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
