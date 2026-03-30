package mocks

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
)

func TestMockSoftwareInstaller_TracksCalls(t *testing.T) {
	m := &MockSoftwareInstaller{
		SoftwareID:        domain.Brave,
		IsInstalledResult: true,
	}

	installed, err := m.IsInstalled()
	if !installed || err != nil {
		t.Errorf("IsInstalled: got %v, %v", installed, err)
	}

	err = m.Install()
	if err != nil {
		t.Errorf("Install: got %v", err)
	}
	if !m.InstallCalled {
		t.Error("InstallCalled should be true")
	}
}

func TestMockSoftwareInstaller_ReturnsError(t *testing.T) {
	m := &MockSoftwareInstaller{
		SoftwareID: domain.Firefox,
		InstallErr: errors.New("install failed"),
	}
	err := m.Install()
	if err == nil {
		t.Fatal("expected error")
	}
}
