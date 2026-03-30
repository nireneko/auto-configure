package browsers_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/browsers"
	"github.com/so-install/pkg/mocks"
)

func TestFirefoxInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("/usr/bin/firefox", "", nil)
	inst := browsers.NewFirefoxInstaller(m)
	got, err := inst.IsInstalled()
	if err != nil || !got {
		t.Errorf("expected true/nil, got %v/%v", got, err)
	}
}

func TestFirefoxInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found"))
	inst := browsers.NewFirefoxInstaller(m)
	got, err := inst.IsInstalled()
	if err != nil || got {
		t.Errorf("expected false/nil, got %v/%v", got, err)
	}
}

func TestFirefoxInstaller_Install_HappyPath(t *testing.T) {
	m := &mocks.MockExecutor{}
	// 5 steps: wget GPG, tee sources, tee preferences, apt update, apt install
	for i := 0; i < 5; i++ {
		m.AddResponse("", "", nil)
	}
	inst := browsers.NewFirefoxInstaller(m)
	err := inst.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 5 {
		t.Fatalf("expected 5 calls, got %d: %v", len(m.Calls), m.Calls)
	}
	assertCall(t, m.Calls[0], "wget", "-q", "https://packages.mozilla.org/apt/repo-signing-key.gpg", "-O", "/etc/apt/keyrings/packages.mozilla.org.gpg")
	assertCall(t, m.Calls[3], "apt", "update")
	assertCall(t, m.Calls[4], "apt", "install", "-y", "firefox")
}

func TestFirefoxInstaller_Install_AptLockOnFinalStep(t *testing.T) {
	m := &mocks.MockExecutor{}
	for i := 0; i < 4; i++ {
		m.AddResponse("", "", nil)
	}
	m.AddResponse("", "E: Could not get lock /var/lib/dpkg/lock", errors.New("exit 100"))
	inst := browsers.NewFirefoxInstaller(m)
	err := inst.Install()
	var aptErr domain.AptLockError
	if !errors.As(err, &aptErr) {
		t.Errorf("expected AptLockError, got %T: %v", err, err)
	}
}
