package tui

import "github.com/so-install/internal/core/domain"

// OSDetectedMsg is sent when OS detection completes.
type OSDetectedMsg struct {
	Info *domain.OSInfo
}

// InstallProgressMsg reports a single software's installation result.
type InstallProgressMsg struct {
	Result domain.InstallResult
}

// StepStartedMsg is sent when a new installation phase starts.
type StepStartedMsg struct {
	Step domain.InstallStep
}

// StepFinishedMsg is sent when an installation phase completes.
type StepFinishedMsg struct {
	Step    domain.InstallStep
	Results []domain.InstallResult
}

// AllInstallsDoneMsg is sent when all selected software has been processed.
type AllInstallsDoneMsg struct {
	Results []domain.InstallResult
}
