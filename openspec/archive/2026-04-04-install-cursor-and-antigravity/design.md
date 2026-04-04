# Design: Add Cursor and Google Antigravity IDEs

## Technical Approach
The implementation will follow the existing patterns for software installation in this project.
- New `SoftwareID` values for `cursor` and `antigravity`.
- New packages in `internal/infrastructure/` for each installer.
- Update `main.go` to register the new installers.

## Architecture Decisions

### Decision: Installer Package Structure
**Choice**: Separate packages `internal/infrastructure/cursor` and `internal/infrastructure/antigravity`.
**Alternatives considered**: A single `internal/infrastructure/ides` package.
**Rationale**: Existing structure uses one package per software (e.g., `docker`, `vscode`, `nvm`), except for `browsers`. To stay consistent with most of the codebase, separate packages are preferred.

### Decision: Antigravity Installation Method
**Choice**: Apt Repository.
**Alternatives considered**: Manual download of `.tar.gz`.
**Rationale**: Using the repository ensures the IDE and CLI are correctly integrated into the system and receive updates via `apt`, which is the primary package manager handled by this tool.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `Cursor` and `Antigravity` constants, update `GetSteps` and `AllSoftware`. |
| `internal/infrastructure/cursor/cursor.go` | Create | Implementation of `CursorInstaller`. |
| `internal/infrastructure/cursor/cursor_test.go` | Create | Unit tests for `CursorInstaller`. |
| `internal/infrastructure/antigravity/antigravity.go` | Create | Implementation of `AntigravityInstaller`. |
| `internal/infrastructure/antigravity/antigravity_test.go` | Create | Unit tests for `AntigravityInstaller`. |
| `cmd/so-install/main.go` | Modify | Register new installers in the `installerMap`. |

## Interfaces / Contracts

The new installers will implement the `domain.SoftwareInstaller` interface:
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
| Unit | `IsInstalled` | Mock `Executor` to return success/failure for `which` commands. |
| Unit | `Install` | Mock `Executor` to verify the sequence of `wget`, `apt`, `gpg`, and `tee` commands. |

## Migration / Rollout
No migration required. The new software options will appear in the TUI next time the application is run.
