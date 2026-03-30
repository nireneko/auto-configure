# Design: Bitwarden via Flatpak

## Technical Approach
We will implement a generic `FlatpakAppInstaller` that satisfies the `domain.SoftwareInstaller` interface. This installer will encapsulate the logic for installing and checking the status of any application available on Flathub. We will then use this generic installer to provide Bitwarden support.

## Architecture Decisions

### Decision: Generic Flatpak App Installer
**Choice**: Create `FlatpakAppInstaller` struct that takes `appID` (e.g., `com.bitwarden.desktop`).
**Alternatives considered**: Create a specific `BitwardenInstaller`.
**Rationale**: Consistency with the `NpmInstaller` pattern and better scalability for future Flatpak applications.

### Decision: Status Check via `flatpak info`
**Choice**: Use `flatpak info <appID>` to check if an app is installed.
**Alternatives considered**: Check for specific files in `/var/lib/flatpak/app/` or `~/.local/share/flatpak/app/`.
**Rationale**: Using the official CLI is more robust and handles both system and user installations correctly.

## Data Flow
The `InstallSoftwareUseCase` calls `IsInstalled()` and then `Install()` on the `FlatpakAppInstaller`. The installer uses the `domain.Executor` to run shell commands.

    InstallSoftwareUseCase ──→ FlatpakAppInstaller ──→ ShellExecutor ──→ flatpak CLI

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `Bitwarden` constant and update `GetSteps()` and `DisplayName()`. |
| `internal/infrastructure/flatpak/app_installer.go` | Create | Implementation of the generic Flatpak app installer. |
| `internal/infrastructure/flatpak/app_installer_test.go` | Create | Unit tests for the new installer. |
| `cmd/so-install/main.go` | Modify | Instantiate and register the Bitwarden installer in the map. |

## Interfaces / Contracts

```go
type FlatpakAppInstaller struct {
	executor    domain.Executor
	appID       string
	softwareID  domain.SoftwareID
}

func NewFlatpakAppInstaller(executor domain.Executor, appID string, id domain.SoftwareID) *FlatpakAppInstaller
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `IsInstalled` | Mock `Executor` to return success/error for `flatpak info`. |
| Unit | `Install` | Mock `Executor` to verify `flatpak install -y flathub ...` call. |
| Integration | End-to-end (manual) | Run the TUI and verify Bitwarden installation on a real system. |

## Migration / Rollout
No migration required. This is a new feature.

## Open Questions
None.
