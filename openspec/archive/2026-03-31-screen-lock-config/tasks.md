# Tasks: Screen Lock Auto-configuration

## Phase 1: Foundation / Domain
- [x] 1.1 Add `ScreenLockConfig` to `SoftwareID` in `internal/core/domain/software.go`.
- [x] 1.2 Add `ScreenLockConfig` to `GetSteps()` and `AllSoftware()` in `internal/core/domain/software.go`.
- [x] 1.3 Update `DisplayName()` for `ScreenLockConfig` in `internal/core/domain/software.go`.

## Phase 2: Infrastructure Implementation
- [x] 2.1 Create `internal/infrastructure/desktop/screen_lock.go` with `ScreenLockInstaller`.
- [x] 2.2 Implement `Install()` with GNOME (`gsettings`) logic.
- [x] 2.3 Implement `Install()` with KDE (`kwriteconfig` + DBus) logic.
- [x] 2.4 Handle `sudo -u $SUDO_USER` wrapping for DE commands.
- [x] 2.5 Implement `IsInstalled()` by reading current DE settings (idempotency).

## Phase 3: Integration / Wiring
- [x] 3.1 Register `ScreenLockInstaller` in `cmd/so-install/main.go`.

## Phase 4: Testing
- [x] 4.1 Create `internal/infrastructure/desktop/screen_lock_test.go`.
- [x] 4.2 Test GNOME/KDE command generation and execution.
- [x] 4.3 Test `sudo` context wrapping.
- [x] 4.4 Test idempotency in `IsInstalled()`.
