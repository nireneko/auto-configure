package desktop

import (
	"github.com/so-install/internal/core/domain"
)

type ScreenLockInstaller struct {
	executor   domain.Executor
	osDetector domain.OSDetector
	userFn     func() string
}

func NewScreenLockInstaller(executor domain.Executor, osDetector domain.OSDetector) *ScreenLockInstaller {
	return &ScreenLockInstaller{
		executor:   executor,
		osDetector: osDetector,
		userFn:     domain.GetActualUser,
	}
}

func (i *ScreenLockInstaller) SetUserFn(fn func() string) {
	i.userFn = fn
}

func (i *ScreenLockInstaller) ID() domain.SoftwareID {
	return domain.ScreenLockConfig
}

func (i *ScreenLockInstaller) IsInstalled() (bool, error) {
	info, err := i.osDetector.Detect()
	if err != nil {
		return false, err
	}

	if info.DesktopEnvironment == domain.GNOME {
		return i.isInstalledGnome()
	}

	if info.DesktopEnvironment == domain.KDE {
		return i.isInstalledKDE()
	}

	return false, nil
}

func (i *ScreenLockInstaller) isInstalledGnome() (bool, error) {
	checks := []struct {
		key   string
		path  string
		value string
	}{
		{"idle-delay", "org.gnome.desktop.session", "uint32 900"},
		{"lock-delay", "org.gnome.desktop.screensaver", "uint32 15"},
		{"lock-enabled", "org.gnome.desktop.screensaver", "true"},
	}

	for _, check := range checks {
		name, args := i.WrapUserCommand("gsettings", []string{"get", check.path, check.key})
		stdout, _, err := i.executor.Execute(name, args...)
		if err != nil {
			return false, nil // Assume not installed if command fails
		}
		if stdout != check.value {
			return false, nil
		}
	}
	return true, nil
}

func (i *ScreenLockInstaller) isInstalledKDE() (bool, error) {
	kread := "kreadconfig5"
	// Check if kwriteconfig6/kreadconfig6 is available
	if i.commandExists("kwriteconfig6") {
		kread = "kreadconfig6"
	}

	checks := []struct {
		key   string
		value string
	}{
		{"Timeout", "900"},
		{"LockGrace", "15"},
		{"Lock", "true"},
	}

	for _, check := range checks {
		name, args := i.WrapUserCommand(kread, []string{"--file", "kscreenlockerrc", "--group", "Daemon", "--key", check.key})
		stdout, _, err := i.executor.Execute(name, args...)
		if err != nil {
			return false, nil
		}
		if stdout != check.value {
			return false, nil
		}
	}
	return true, nil
}

func (i *ScreenLockInstaller) Install() error {
	info, err := i.osDetector.Detect()
	if err != nil {
		return err
	}

	if info.DesktopEnvironment == domain.GNOME {
		return i.installGnome()
	}

	if info.DesktopEnvironment == domain.KDE {
		return i.installKDE()
	}

	return nil
}

func (i *ScreenLockInstaller) installKDE() error {
	kwrite := "kwriteconfig5"
	// Check if kwriteconfig6 is available
	if i.commandExists("kwriteconfig6") {
		kwrite = "kwriteconfig6"
	}

	commands := [][]string{
		{kwrite, "--file", "kscreenlockerrc", "--group", "Daemon", "--key", "Timeout", "900"},
		{kwrite, "--file", "kscreenlockerrc", "--group", "Daemon", "--key", "LockGrace", "15"},
		{kwrite, "--file", "kscreenlockerrc", "--group", "Daemon", "--key", "Lock", "true"},
		{"dbus-send", "--session", "--dest=org.kde.kscreenlocker", "--type=method_call", "/Main", "org.kde.kscreenlocker.configure"},
	}

	for _, cmd := range commands {
		name, args := i.WrapUserCommand(cmd[0], cmd[1:])
		_, _, err := i.executor.Execute(name, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ScreenLockInstaller) commandExists(name string) bool {
	_, _, err := i.executor.Execute("which", name)
	return err == nil
}

func (i *ScreenLockInstaller) installGnome() error {
	commands := [][]string{
		{"gsettings", "set", "org.gnome.desktop.session", "idle-delay", "900"},
		{"gsettings", "set", "org.gnome.desktop.screensaver", "lock-delay", "15"},
		{"gsettings", "set", "org.gnome.desktop.screensaver", "lock-enabled", "true"},
	}

	for _, cmd := range commands {
		name, args := i.WrapUserCommand(cmd[0], cmd[1:])
		_, _, err := i.executor.Execute(name, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ScreenLockInstaller) WrapUserCommand(name string, args []string) (string, []string) {
	actualUser := i.userFn()
	if actualUser != "" && actualUser != "root" {
		newArgs := append([]string{"-u", actualUser, name}, args...)
		return "sudo", newArgs
	}
	return name, args
}
