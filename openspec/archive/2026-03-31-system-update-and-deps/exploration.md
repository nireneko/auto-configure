## Exploration: System Update and Base Dependencies

### Current State
The system orchestrates software installations through `InstallSoftwareUseCase`, which follows steps defined in `domain.GetSteps()`. The selection is made in the TUI, which picks from `domain.AllSoftware()`. Currently, there is no mandatory pre-installation step for system updates or base dependencies.

### Affected Areas
- `internal/core/domain/software.go` — To add `SystemUpdate` and `BaseDeps` to `SoftwareID` and `GetSteps()`.
- `internal/infrastructure/apt/` (New) — To implement `SystemUpdateInstaller` and `BaseDepsInstaller`.
- `cmd/so-install/main.go` — To register the new installers.
- `internal/presentation/tui/model.go` — To ensure these mandatory steps are included in the execution even if not selectable in the UI.

### Approaches
1. **Mandatory Software IDs** — Add `SystemUpdate` and `BaseDeps` as `SoftwareID`s. Add them to the first step in `GetSteps()`. Modify the TUI to always prepend these to the selected list when starting installation.
   - Pros: Consistent with existing architecture, visible progress in TUI.
   - Cons: Requires minor modification to TUI selection logic.
   - Effort: Low

2. **Hardcoded Hook in UseCase** — Modify `InstallSoftwareUseCase.Execute` to run `apt update` and `apt upgrade` before the loop.
   - Pros: Simple to implement.
   - Cons: Not visible in TUI progress, bypasses the "step-based" pattern, harder to test in isolation.
   - Effort: Low

3. **Pre-Selection State** — Add them to `AllSoftware()` and check them by default in the UI.
   - Pros: Very transparent.
   - Cons: User could uncheck them (violating "ensure they run first").
   - Effort: Low

### Recommendation
**Approach 1** is recommended. It keeps the architecture clean by reusing the `SoftwareInstaller` interface and `InstallStep` logic, while ensuring the steps are visible to the user during progress but not bypassable in the selection phase.

### Risks
- **`apt upgrade` duration**: Upgrading the whole system can take a long time, potentially frustrating users who just want one app.
- **Interactivity**: `apt upgrade` can sometimes prompt for configuration file changes. We should use `DEBIAN_FRONTEND=noninteractive`.
- **Apt Lock**: Running `apt update` at the start might still hit the apt lock if the system is doing background updates. The use case already has retry logic for this.

### Ready for Proposal
Yes. The orchestrator should proceed with the proposal to add these mandatory system preparation steps.
