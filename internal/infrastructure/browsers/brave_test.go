package browsers_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/browsers"
	"github.com/so-install/pkg/mocks"
)

func TestBraveInstaller_IsInstalled_True(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("/usr/bin/brave-browser", "", nil)
	inst := browsers.NewBraveInstaller(m)
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected true when which succeeds")
	}
	if len(m.Calls) != 1 || m.Calls[0].Name != "which" {
		t.Errorf("expected 'which' call, got %v", m.Calls)
	}
}

func TestBraveInstaller_IsInstalled_False(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", errors.New("not found"))
	inst := browsers.NewBraveInstaller(m)
	got, err := inst.IsInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected false when which fails")
	}
}

func TestBraveInstaller_Install_HappyPath(t *testing.T) {
	m := &mocks.MockExecutor{}
	// 5 commands: mkdir, wget gpg, sh -c tee repo, apt update, apt install
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	m.AddResponse("", "", nil)
	inst := browsers.NewBraveInstaller(m)
	err := inst.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Calls) != 5 {
		t.Fatalf("expected 5 calls, got %d: %v", len(m.Calls), m.Calls)
	}
	// Verify exact command sequence
	assertCall(t, m.Calls[0], "mkdir", "-p", "/usr/share/keyrings")
	assertCall(t, m.Calls[1], "wget", "-qO", "/usr/share/keyrings/brave-browser-archive-keyring.gpg", "https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg")
	assertCall(t, m.Calls[2], "sh", "-c", "echo 'deb [signed-by=/usr/share/keyrings/brave-browser-archive-keyring.gpg] https://brave-browser-apt-release.s3.brave.com/ stable main' | tee /etc/apt/sources.list.d/brave-browser-release.list")
	assertCall(t, m.Calls[3], "apt", "update")
	assertCall(t, m.Calls[4], "apt", "install", "-y", "brave-browser")
}

func TestBraveInstaller_Install_StopsOnFirstFailure(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "permission denied", errors.New("exit 1"))
	inst := browsers.NewBraveInstaller(m)
	err := inst.Install()
	if err == nil {
		t.Fatal("expected error on failure")
	}
	if len(m.Calls) != 1 {
		t.Errorf("expected only 1 call before stopping, got %d", len(m.Calls))
	}
}

func TestBraveInstaller_Install_AptLockReturnsAptLockError(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil) // mkdir ok
	m.AddResponse("", "", nil) // wget gpg ok
	m.AddResponse("", "", nil) // sh tee repo ok
	m.AddResponse("", "", nil) // apt update ok
	m.AddResponse("", "E: Could not get lock /var/lib/dpkg/lock", errors.New("exit 100"))
	inst := browsers.NewBraveInstaller(m)
	err := inst.Install()
	if err == nil {
		t.Fatal("expected error")
	}
	var aptErr domain.AptLockError
	if !errors.As(err, &aptErr) {
		t.Errorf("expected AptLockError, got %T: %v", err, err)
	}
}

func TestBraveInstaller_Install_NonLockErrorReturnsInstallError(t *testing.T) {
	m := &mocks.MockExecutor{}
	m.AddResponse("", "", nil) // mkdir ok
	m.AddResponse("", "", nil) // wget gpg ok
	m.AddResponse("", "", nil) // sh tee repo ok
	m.AddResponse("", "", nil) // apt update ok
	m.AddResponse("", "E: Package not found", errors.New("exit 100"))
	inst := browsers.NewBraveInstaller(m)
	err := inst.Install()
	var aptErr domain.AptLockError
	if errors.As(err, &aptErr) {
		t.Error("expected InstallError not AptLockError")
	}
	var instErr domain.InstallError
	if !errors.As(err, &instErr) {
		t.Errorf("expected InstallError, got %T", err)
	}
}

// assertCall is a helper to verify a specific executor call.
func assertCall(t *testing.T, call mocks.ExecutorCall, name string, args ...string) {
	t.Helper()
	if call.Name != name {
		t.Errorf("expected command %q, got %q", name, call.Name)
	}
	if len(call.Args) != len(args) {
		t.Errorf("command %q: expected args %v, got %v", name, args, call.Args)
		return
	}
	for i, a := range args {
		if call.Args[i] != a {
			t.Errorf("command %q arg[%d]: expected %q, got %q", name, i, a, call.Args[i])
		}
	}
}
