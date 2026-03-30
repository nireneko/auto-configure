# Proposal: System Update and Base Dependencies

## Intent
Ensure the system is up-to-date (`apt update && apt upgrade`) and has essential tools (`git`, `wget`, `curl`, etc.) installed before any other software installation. This provides a stable and consistent base for subsequent steps.

## Scope

### In Scope
- Add `SystemUpdate` and `BaseDeps` to `domain.SoftwareID`.
- Add a "System Preparation" step at the beginning of `domain.GetSteps()`.
- Implement `AptUpdateInstaller` (runs `apt update` and `apt upgrade -y`).
- Implement `BaseDepsInstaller` (installs `git`, `wget`, `curl`, `ca-certificates`, `gnupg`, `lsb-release`).
- Modify TUI to always include these steps in the installation sequence.

### Out of Scope
- Making these steps optional (user requested they run first).
- Installing software-specific dependencies (these should remain in their respective installers).

## Approach
Reusing the `SoftwareInstaller` interface, we will create two new implementations in `internal/infrastructure/apt`. These will be registered in `main.go`. The TUI will be updated to automatically select these IDs when the user confirms their selection, ensuring they are executed in the first step.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Add new IDs and the first step in `GetSteps()`. |
| `internal/infrastructure/apt/` | New | Create `update.go` and `deps.go` with installers. |
| `cmd/so-install/main.go` | Modified | Register new installers in the map. |
| `internal/presentation/tui/model.go` | Modified | Prepend mandatory IDs to selection on "enter". |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Long upgrade time | High | Visible progress in TUI; mandatory as per request. |
| Upgrade interactivity | Medium | Use `DEBIAN_FRONTEND=noninteractive`. |
| APT lock | Medium | Existing retry logic in `InstallSoftwareUseCase`. |

## Rollback Plan
- Revert `domain.GetSteps()` to remove the first step.
- Remove `AptUpdate` and `BaseDeps` from the TUI selection logic.

## Success Criteria
- [ ] `apt update` and `apt upgrade` run successfully at start.
- [ ] `git`, `wget`, `curl` are installed if not present.
- [ ] Installation proceeds to the next steps.
