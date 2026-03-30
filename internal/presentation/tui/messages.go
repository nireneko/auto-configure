package tui

import "github.com/so-install/internal/core/domain"

// OSDetectedMsg is sent when OS detection completes.
type OSDetectedMsg struct {
	Info *domain.OSInfo
}

// InstallProgressMsg reports a single browser's installation result.
type InstallProgressMsg struct {
	Result domain.InstallResult
}

// AllInstallsDoneMsg is sent when all selected browsers have been processed.
type AllInstallsDoneMsg struct {
	Results []domain.InstallResult
}
