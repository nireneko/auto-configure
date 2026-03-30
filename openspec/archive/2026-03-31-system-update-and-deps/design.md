# Design: System Update and Base Dependencies

## Technical Approach
Implement two new `domain.SoftwareInstaller`s: `AptUpdateInstaller` and `BaseDepsInstaller`. These will be registered in `main.go` and added to a new "system-prep" step in `domain.GetSteps()`. The TUI will be modified to ensure these IDs are always included in the selection list when starting the installation process.

## Architecture Decisions

### Decision: Mandatory Step Implementation

**Choice**: Add `SystemUpdate` and `BaseDeps` to `domain.GetSteps()` and prepend them to the selected list in the TUI.
**Alternatives considered**: Hardcoding the steps in `InstallSoftwareUseCase.Execute`.
**Rationale**: Reusing the existing step-based orchestration allows for better visibility in the TUI progress view and keeps the use case generic.

### Decision: System Update Strategy

**Choice**: Run both `apt-get update` and `apt-get upgrade -y`.
**Alternatives considered**: Only `apt-get update`.
**Rationale**: The user explicitly requested both to ensure the system is fully up-to-date before installing new software.

## Data Flow
The TUI model will receive the user's selection, then prepend the mandatory system preparation IDs before passing the list to the installation command.

```
TUI (Selection) ──→ Model.handleKey("enter") ──→ Prepend [SystemUpdate, BaseDeps] ──→ TUI (Progress) ──→ UseCase.Execute
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `SystemUpdate`, `BaseDeps` and `system-prep` step. |
| `internal/infrastructure/apt/update.go` | Create | `AptUpdateInstaller` implementation. |
| `internal/infrastructure/apt/deps.go` | Create | `BaseDepsInstaller` implementation. |
| `cmd/so-install/main.go` | Modify | Register new installers in the map. |
| `internal/presentation/tui/model.go` | Modify | Prepend mandatory IDs on "enter" key. |

## Interfaces / Contracts

```go
// internal/infrastructure/apt/update.go
type AptUpdateInstaller struct {
    executor domain.Executor
}

// internal/infrastructure/apt/deps.go
type BaseDepsInstaller struct {
    executor domain.Executor
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `AptUpdateInstaller` | Mock executor to verify `apt-get update` and `apt-get upgrade` calls. |
| Unit | `BaseDepsInstaller` | Mock executor to verify `apt-get install` with the correct packages. |
| Unit | `domain.GetSteps` | Verify the "system-prep" step is first and has the correct IDs. |
| Unit | `TUI Model` | Verify mandatory IDs are prepended to selection. |

## Migration / Rollout
No migration required. This is a new feature that will run for all future installations.

## Open Questions
- None.
