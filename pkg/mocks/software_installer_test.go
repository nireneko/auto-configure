package mocks

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
)

func TestMockSoftwareInstaller(t *testing.T) {
	m := &MockSoftwareInstaller{
		SoftwareID:        domain.Brave,
		IsInstalledResult: true,
		InstallErr:        errors.New("fail"),
	}

	if m.ID() != domain.Brave {
		t.Errorf("expected ID Brave, got %v", m.ID())
	}

	installed, err := m.IsInstalled()
	if err != nil {
		t.Errorf("unexpected IsInstalled error: %v", err)
	}
	if !installed {
		t.Error("expected IsInstalled to be true")
	}

	err = m.Install()
	if err == nil {
		t.Error("expected Install error, got nil")
	}
	if !m.InstallCalled {
		t.Error("expected InstallCalled to be true")
	}
}
