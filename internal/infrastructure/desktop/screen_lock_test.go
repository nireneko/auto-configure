package desktop_test

import (
	"errors"
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/desktop"
	"github.com/so-install/pkg/mocks"
)

func TestScreenLockInstaller_ID(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	installer := desktop.NewScreenLockInstaller(executor, osDetector)

	if installer.ID() != domain.ScreenLockConfig {
		t.Errorf("Expected ID %s, got %s", domain.ScreenLockConfig, installer.ID())
	}
}

func TestScreenLockInstaller_Install_GNOME(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.GNOME
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	// Override userFn to avoid sudo wrapping in test
	installer.SetUserFn(func() string { return "root" })

	err := installer.Install()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(executor.Calls) != 3 {
		t.Fatalf("Expected 3 calls, got %d", len(executor.Calls))
	}

	expected := []struct {
		name string
		args []string
	}{
		{"gsettings", []string{"set", "org.gnome.desktop.session", "idle-delay", "900"}},
		{"gsettings", []string{"set", "org.gnome.desktop.screensaver", "lock-delay", "15"}},
		{"gsettings", []string{"set", "org.gnome.desktop.screensaver", "lock-enabled", "true"}},
	}

	for i, call := range executor.Calls {
		if call.Name != expected[i].name {
			t.Errorf("Call %d Name = %s, want %s", i, call.Name, expected[i].name)
		}
		for j, arg := range call.Args {
			if arg != expected[i].args[j] {
				t.Errorf("Call %d Arg %d = %s, want %s", i, j, arg, expected[i].args[j])
			}
		}
	}
}

func TestScreenLockInstaller_Install_KDE(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.KDE
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })

	// First call is commandExists("kwriteconfig6")
	executor.AddResponse("", "", errors.New("not found")) // Not found, fallback to kwriteconfig5

	err := installer.Install()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(executor.Calls) != 5 {
		t.Fatalf("Expected 5 calls (1 check + 4 configs), got %d", len(executor.Calls))
	}

	// First call was the check
	if executor.Calls[0].Name != "which" {
		t.Errorf("First call Name = %s, want which", executor.Calls[0].Name)
	}

	expected := []struct {
		name string
		args []string
	}{
		{"kwriteconfig5", []string{"--file", "kscreenlockerrc", "--group", "Daemon", "--key", "Timeout", "900"}},
		{"kwriteconfig5", []string{"--file", "kscreenlockerrc", "--group", "Daemon", "--key", "LockGrace", "15"}},
		{"kwriteconfig5", []string{"--file", "kscreenlockerrc", "--group", "Daemon", "--key", "Lock", "true"}},
		{"dbus-send", []string{"--session", "--dest=org.kde.kscreenlocker", "--type=method_call", "/Main", "org.kde.kscreenlocker.configure"}},
	}

	for i, call := range executor.Calls[1:] {
		if call.Name != expected[i].name {
			t.Errorf("Call %d Name = %s, want %s", i+1, call.Name, expected[i].name)
		}
	}
}

func TestScreenLockInstaller_WrapUserCommand(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	installer := desktop.NewScreenLockInstaller(executor, osDetector)

	// Case 1: Running as root, should NOT wrap
	installer.SetUserFn(func() string { return "root" })
	name, args := installer.WrapUserCommand("gsettings", []string{"get", "path", "key"})
	if name != "gsettings" {
		t.Errorf("Expected gsettings, got %s", name)
	}

	// Case 2: Running as user via sudo, SHOULD wrap
	installer.SetUserFn(func() string { return "borja" })
	name, args = installer.WrapUserCommand("gsettings", []string{"get", "path", "key"})
	if name != "sudo" {
		t.Errorf("Expected sudo, got %s", name)
	}
	if args[0] != "-u" || args[1] != "borja" || args[2] != "gsettings" {
		t.Errorf("Incorrect sudo wrapping: %v", args)
	}
}

func TestScreenLockInstaller_IsInstalled_GNOME(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.GNOME
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })

	// Case 1: All settings match
	executor.AddResponse("uint32 900", "", nil) // idle-delay
	executor.AddResponse("uint32 15", "", nil)  // lock-delay
	executor.AddResponse("true", "", nil)       // lock-enabled

	installed, err := installer.IsInstalled()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !installed {
		t.Error("Expected IsInstalled to be true when all settings match")
	}

	// Case 2: One setting mismatch
	executor.AddResponse("uint32 600", "", nil) // idle-delay
	// Other calls won't happen if it returns early

	installed, err = installer.IsInstalled()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if installed {
		t.Error("Expected IsInstalled to be false when settings mismatch")
	}
}

func TestScreenLockInstaller_IsInstalled_KDE(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.KDE
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })

	// Case 1: All settings match
	executor.AddResponse("", "", errors.New("not found")) // which kwriteconfig6
	executor.AddResponse("900", "", nil)                  // Timeout
	executor.AddResponse("15", "", nil)                   // LockGrace
	executor.AddResponse("true", "", nil)                 // Lock

	installed, err := installer.IsInstalled()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !installed {
		t.Error("Expected IsInstalled to be true when all settings match")
	}
}
