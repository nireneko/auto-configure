package desktop_test

import (
	"testing"

	"github.com/so-install/internal/core/domain"
	"github.com/so-install/internal/infrastructure/desktop"
	"github.com/so-install/pkg/mocks"
)

func TestScreenLockInstaller_ID(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	if installer.ID() != domain.ScreenLockConfig { t.Errorf("Expected ID %s, got %s", domain.ScreenLockConfig, installer.ID()) }
}

func TestScreenLockInstaller_Install_GNOME(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.GNOME
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })
	err := installer.Install()
	if err != nil { t.Fatalf("Expected no error, got %v", err) }
}

func TestScreenLockInstaller_Install_KDE(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.KDE
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })
	executor.AddResponse("", "", nil) // which kwrite6 -> ok
	err := installer.Install()
	if err != nil { t.Fatalf("Expected no error, got %v", err) }
}

func TestScreenLockInstaller_IsInstalled_GNOME(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.GNOME
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })
	executor.AddResponse("uint32 900", "", nil)
	executor.AddResponse("uint32 15", "", nil)
	executor.AddResponse("true", "", nil)
	installed, _ := installer.IsInstalled()
	if !installed { t.Error("expected true") }
}

func TestScreenLockInstaller_IsInstalled_KDE_k6(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.KDE
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installer.SetUserFn(func() string { return "root" })
	executor.AddResponse("", "", nil) // which k6 -> ok
	executor.AddResponse("900", "", nil)
	executor.AddResponse("15", "", nil)
	executor.AddResponse("true", "", nil)
	installed, _ := installer.IsInstalled()
	if !installed { t.Error("expected true") }
}

func TestScreenLockInstaller_IsInstalled_Other(t *testing.T) {
	executor := new(mocks.MockExecutor)
	osDetector := new(mocks.MockOSDetector)
	osDetector.ReturnDE = domain.Other
	installer := desktop.NewScreenLockInstaller(executor, osDetector)
	installed, err := installer.IsInstalled()
	if err != nil { t.Error("unexpected error") }
	if installed { t.Error("expected false") }
}
