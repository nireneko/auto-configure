# Design: Install VS Code

## Technical Approach
We will implement the `VsCodeInstaller` in `internal/infrastructure/vscode/vscode.go`, following the established `SoftwareInstaller` interface. The installation will use `wget` to download the official `.deb` from Microsoft and `apt` to install it. This approach is consistent with `ChromeInstaller`.

## Architecture Decisions

### Decision: Installation Method
**Choice**: Official `.deb` package via `apt`.
**Alternatives considered**: Flatpak, Snap, Manual tarball.
**Rationale**: Native `.deb` packages provide better integration with the system (PATH, icons, mime types) and automatically configure the Microsoft repository for updates. It's the standard for development on Debian/Ubuntu.

### Decision: Installer Location
**Choice**: `internal/infrastructure/vscode/`.
**Alternatives considered**: `internal/infrastructure/ides/`.
**Rationale**: The project currently groups infrastructure by specific software or category (e.g., `browsers`, `docker`, `ollama`). Since we only have VS Code for now, a dedicated package is cleaner, but we'll add an `ides` step in the domain.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `VsCode` SoftwareID and update `GetSteps()`/`AllSoftware()`. |
| `internal/core/domain/software_test.go` | Modify | Add display name test for VS Code. |
| `internal/infrastructure/vscode/vscode.go` | Create | Implementation of `VsCodeInstaller`. |
| `internal/infrastructure/vscode/vscode_test.go` | Create | Unit tests with `MockExecutor`. |
| `cmd/so-install/main.go` | Modify | Register `VsCodeInstaller` in the `installerMap`. |
| `PRD.md` | Modify | Add IDE module to requirements. |

## Interfaces / Contracts
No new interfaces. We use the existing `domain.SoftwareInstaller`:

```go
type SoftwareInstaller interface {
	Install() error
	IsInstalled() (bool, error)
	ID() SoftwareID
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `VsCodeInstaller.Install` | Mock `Executor` to verify `wget` and `apt` commands. |
| Unit | `VsCodeInstaller.IsInstalled` | Mock `Executor` to verify `which code` behavior. |
| Unit | `domain.SoftwareID.DisplayName` | Standard table-driven test in `software_test.go`. |
| Unit | `domain.GetSteps` | Verify `ides` step inclusion and position. |

## Migration / Rollout
No migration required. This is a new feature.

## Open Questions
None. The pattern is well-established in the codebase.
