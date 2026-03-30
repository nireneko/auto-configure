package mocks

import "github.com/so-install/internal/core/domain"

// MockBrowserInstaller records Install/IsInstalled calls and returns configured values.
type MockBrowserInstaller struct {
	BrowserID         domain.BrowserID
	InstallErr        error
	InstallCalled     bool
	IsInstalledResult bool
	IsInstalledErr    error
}

var _ domain.BrowserInstaller = (*MockBrowserInstaller)(nil)

func (m *MockBrowserInstaller) Install() error {
	m.InstallCalled = true
	return m.InstallErr
}

func (m *MockBrowserInstaller) IsInstalled() (bool, error) {
	return m.IsInstalledResult, m.IsInstalledErr
}

func (m *MockBrowserInstaller) ID() domain.BrowserID {
	return m.BrowserID
}
