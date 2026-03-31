## Exploration: Gitlab Token Config for Composer and NPM

### Current State
The tool currently manages software installations via `SoftwareInstaller` interface. It doesn't have a way to handle interactive user input during or before the installation process. Configuration files for Composer (`~/.composer/auth.json`) and NPM (`~/.npmrc`) are not managed.

### Affected Areas
- `internal/core/domain/software.go` — Add `GitlabConfig` ID and update `GetSteps` / `AllSoftware`.
- `internal/infrastructure/` — Create `internal/infrastructure/gitlab/configurator.go` to handle the file writing.
- `internal/presentation/tui/model.go` — Add `stateTokenInput` to `appState` and handle the transition.
- `internal/presentation/tui/view.go` — Update views to include the token input field.
- `cmd/so-install/main.go` — Register the new `GitlabConfigurator`.

### Approaches
1. **Extend `SoftwareInstaller` with `Configure` step** — Add a way for the TUI to know if an installer needs configuration before proceeding.
   - Pros: Clean separation of concerns.
   - Cons: More changes to core interfaces.
   - Effort: Medium

2. **Special handling in TUI for `GitlabConfig`** — If `GitlabConfig` is selected, trigger a specific input state.
   - Pros: Easier to implement within the current TUI structure.
   - Cons: Hardcodes logic for one specific software ID.
   - Effort: Low

3. **Treat configuration as "Installation" with a pre-set token** — If we can get the token from an environment variable or a command line flag, we wouldn't need a TUI input.
   - Pros: Simple.
   - Cons: Not user-friendly (the user wants to input it).
   - Effort: Very Low

### Recommendation
Approach 2: **Special handling in TUI**. We'll add a `stateTokenInput` and if `GitlabConfig` is selected, the TUI will ask for the token before starting the installation sequence. This keeps the current `SoftwareInstaller` interface as is but adds the necessary interaction.

### Risks
- Overwriting existing `~/.npmrc` or `~/.composer/auth.json` accidentally. We should append or use `composer config` / `npm config` commands instead of direct file manipulation if possible.
- Security: We must ensure the token is handled as a secret (no echoing it in logs, although TUI input will show it by default unless we use a masked field).

### Ready for Proposal
Yes — I have a clear plan for both the backend logic and the TUI changes.
