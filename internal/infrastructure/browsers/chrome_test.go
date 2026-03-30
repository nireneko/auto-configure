package browsers_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/browsers"
	"github.com/so-install/pkg/mocks"
)

func TestChromeInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("/usr/bin/google-chrome-stable", "", nil)
	inst := browsers.NewChromeInstaller(m)
	got, err := inst.IsInstalled()
	if err != nil || !got {
		t.Errorf("expected true/nil, got %v/%v", got, err)
	}
}

func TestChromeInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found"))
	inst := browsers.NewChromeInstaller(m)
	got, err := inst.IsInstalled()
	if err != nil || got {
		t.Errorf("expected false/nil, got %v/%v", got, err)
	}
}

func TestChromeInstaller_Install_HappyPath(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	inst := browsers.NewChromeInstaller(m)
	err := inst.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(m.Calls))
	}
	assertCall(t, m.Calls[0], "wget", "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb", "-P", "/tmp/")
	assertCall(t, m.Calls[1], "apt", "install", "-y", "/tmp/google-chrome-stable_current_amd64.deb")
}

func TestChromeInstaller_Install_AptLock(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil)
	m.AddResponse("", "E: Could not get lock /var/lib/dpkg/lock", errors.New("exit 100"))
	inst := browsers.NewChromeInstaller(m)
	err := inst.Install()
	var aptErr domain.AptLockError
	if !errors.As(err, &aptErr) {
		t.Errorf("expected AptLockError, got %T", err)
	}
}
