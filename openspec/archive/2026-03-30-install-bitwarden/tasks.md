# Tasks: Bitwarden via Flatpak

## Phase 1: Foundation & Domain Update
- [ ] 1.1 Add `Bitwarden` constant to `SoftwareID` in `internal/core/domain/software.go`.
- [ ] 1.2 Update `GetSteps()` in `internal/core/domain/software.go` to include Bitwarden in a new "apps" step after Flatpak.
- [ ] 1.3 Update `DisplayName()` in `internal/core/domain/software.go` for `Bitwarden`.
- [ ] 1.4 Add `Bitwarden` to `AllSoftware()` in `internal/core/domain/software.go`.

## Phase 2: Infrastructure Implementation
- [ ] 2.1 Create `internal/infrastructure/flatpak/app_installer.go` with `FlatpakAppInstaller` struct.
- [ ] 2.2 Implement `ID()` method for `FlatpakAppInstaller`.
- [ ] 2.3 Implement `IsInstalled()` using `flatpak info <appID>`.
- [ ] 2.4 Implement `Install()` using `flatpak install -y flathub <appID>`.

## Phase 3: Wiring & Main
- [ ] 3.1 Register `Bitwarden` in the `installerMap` in `cmd/so-install/main.go` using `NewFlatpakAppInstaller`.

## Phase 4: Testing & Verification
- [ ] 4.1 Create unit tests in `internal/infrastructure/flatpak/app_installer_test.go`.
- [ ] 4.2 Test `IsInstalled` happy path (app exists).
- [ ] 4.3 Test `IsInstalled` unhappy path (app doesn't exist).
- [ ] 4.4 Test `Install` execution and error handling.
- [ ] 4.5 Verify TUI shows Bitwarden and triggers installation.
