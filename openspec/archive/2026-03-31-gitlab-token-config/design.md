# Design: Gitlab Token Configuration

## Technical Approach

The implementation will add a new `GitlabTokenConfigurator` in `internal/infrastructure/gitlab`. This configurator will implement the `SoftwareInstaller` interface, but its "Install" step will actually perform configuration.

To handle the interactive token input, the TUI `Model` will be extended with a new state `stateTokenInput`. If `GitlabTokenConfig` is selected, the model will transition to this state before starting the installation sequence.

## Architecture Decisions

### Decision: Reuse `SoftwareInstaller` Interface

**Choice**: Implement `GitlabTokenConfigurator` as a `SoftwareInstaller`.
**Alternatives considered**: Create a new `Configurator` interface.
**Rationale**: Reusing the existing interface allows us to plug the configuration step directly into the existing `InstallBySteps` and `InstallSoftwareUseCase` without significant refactoring. The "installation" in this case is the act of configuring the system.

### Decision: State-based Token Capture

**Choice**: Add `stateTokenInput` to the TUI model.
**Alternatives considered**: Capture the token during the "installation" phase via stdin.
**Rationale**: Bubbletea is an event-driven framework; it's much cleaner to have a dedicated state and view for input rather than trying to hijack the terminal during an installation step.

### Decision: Manual File Updates vs CLI Tools

**Choice**: Use direct file updates (atomic writes) for `~/.composer/auth.json` and `~/.npmrc`.
**Alternatives considered**: Use `composer config` and `npm config` commands.
**Rationale**: Using the CLI tools requires them to be installed first. Since this configuration might happen *before* or *during* the installation of related tools, direct file manipulation is more robust and doesn't depend on the tools being in the PATH. We will use `domain.GetActualHome()` to ensure we target the correct user's directory.

## Data Flow

1. **Selection**: User selects `GitlabTokenConfig` in `stateSoftwareSelect`.
2. **Transition**: On "Enter", if `GitlabTokenConfig` is selected, state moves to `stateTokenInput`.
3. **Capture**: User enters token (masked).
4. **Injection**: Token is stored in the `Model` and passed to the `GitlabTokenConfigurator` via a setter or constructor update before the installation starts.
5. **Execution**: During the `apps` phase, `GitlabTokenConfigurator.Install()` is called, which writes the configuration files.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `GitlabTokenConfig` ID and update `GetSteps`. |
| `internal/infrastructure/gitlab/configurator.go` | Create | Implementation of `GitlabTokenConfigurator`. |
| `internal/presentation/tui/model.go` | Modify | Add `stateTokenInput`, `gitlabToken` string, and state transition logic. |
| `internal/presentation/tui/view_token.go` | Create | (New file or added to model.go) View for token input. |
| `cmd/so-install/main.go` | Modify | Instantiate and register `GitlabTokenConfigurator`. |

## Interfaces / Contracts

The `GitlabTokenConfigurator` will need a way to receive the token:

```go
type GitlabTokenConfigurator struct {
    token string
    // ... other fields
}

func (g *GitlabTokenConfigurator) SetToken(token string) {
    g.token = token
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `GitlabTokenConfigurator` | Mock `Executor` and `fs` (if possible) to verify correct file content is "written". |
| Unit | TUI State Transition | Test that selecting Gitlab config leads to `stateTokenInput`. |
| Integration | End-to-end file write | Run the configurator in a temporary directory and verify `auth.json` and `.npmrc` content. |

## Migration / Rollout

No migration required. This is a new feature.

## Open Questions

- [ ] Should we support multiple Gitlab domains? (Defaulting to gitlab.com as requested).
- [ ] Should the token be validated in any way (e.g., non-empty, basic format)?
