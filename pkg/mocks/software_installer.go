package mocks

import "github.com/so-install/internal/core/domain"

// MockSoftwareInstaller records Install/IsInstalled calls and returns configured values.
type MockSoftwareInstaller struct {
	SoftwareID        domain.SoftwareID
	InstallErr        error
	InstallCalled     bool
	IsInstalledResult bool
	IsInstalledErr    error
}

var _ domain.SoftwareInstaller = (*MockSoftwareInstaller)(nil)

func (m *MockSoftwareInstaller) Install() error {
	m.InstallCalled = true
	return m.InstallErr
}

func (m *MockSoftwareInstaller) IsInstalled() (bool, error) {
	return m.IsInstalledResult, m.IsInstalledErr
}

func (m *MockSoftwareInstaller) ID() domain.SoftwareID {
	return m.SoftwareID
}
