# Design: Screen Lock Auto-configuration

## Technical Approach
Implement a new infrastructure component `ScreenLockInstaller` that fulfills the `domain.SoftwareInstaller` interface. This installer will detect the Desktop Environment using `domain.OSDetector` and execute the appropriate shell commands via `domain.Executor`.

## Architecture Decisions

### Decision: SoftwareInstaller Interface
**Choice**: Use `domain.SoftwareInstaller`.
**Alternatives considered**: Create a separate `SystemConfigurator` interface.
**Rationale**: The `SoftwareInstaller` interface already provides `Install()`, `IsInstalled()`, and `ID()`, which perfectly fit this use case. Adding it as a "software" allows it to be part of the standard installation steps.

### Decision: User Context Execution
**Choice**: Wrap commands with `sudo -u $USER` when running as root.
**Alternatives considered**: Run directly as root.
**Rationale**: `gsettings` and `kwriteconfig` modify user-specific configuration files and sessions. Running as root would either fail or modify the root user's desktop settings.

## Data Flow
`InstallSoftwareUseCase` → `ScreenLockInstaller.Install()`
  → `OSDetector.Detect()` (Check DE)
  → `Executor.Execute()` (Run commands)

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `ScreenLockConfig` SoftwareID. |
| `internal/infrastructure/desktop/screen_lock.go` | Create | New installer for screen lock configuration. |
| `internal/infrastructure/desktop/screen_lock_test.go` | Create | Unit tests for the installer. |
| `cmd/so-install/main.go` | Modify | Register the new installer. |

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `ScreenLockInstaller` | Mock `Executor` and `OSDetector`. Verify correct commands for GNOME and KDE. |
| Unit | `domain.GetSteps` | Verify `ScreenLockConfig` is included in the steps. |
