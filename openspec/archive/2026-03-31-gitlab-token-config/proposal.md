# Proposal: Gitlab Token Configuration

## Intent
Configure a global Gitlab token for Composer and NPM to allow the user to access private repositories and packages.

## Scope

### In Scope
- New `SoftwareID`: `GitlabConfig`.
- Implementation of `GitlabConfigurator` to update `~/.composer/auth.json` and `~/.npmrc`.
- TUI update with a new `stateTokenInput` to capture the token securely.
- Integration of the configuration step into the installation workflow.

### Out of Scope
- Support for multiple Gitlab domains (defaulting to `gitlab.com`).
- Validation of the token against Gitlab API.

## Approach
- Add `GitlabConfig` to `internal/core/domain/software.go`.
- If `GitlabConfig` is selected, the TUI transitions to `stateTokenInput` after the selection phase.
- The `GitlabConfigurator` will read/update the configuration files in the user's home directory (handling `sudo` correctly via `GetActualHome`).
- Use `composer config` and `npm config` commands if available, otherwise direct file manipulation with safe-writing patterns.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Add `GitlabConfig` ID and update steps. |
| `internal/infrastructure/gitlab/` | New | Create `configurator.go` for the logic. |
| `internal/presentation/tui/model.go` | Modified | Add `stateTokenInput` and input handling logic. |
| `cmd/so-install/main.go` | Modified | Register the new configurator. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Token exposure in logs | Low | Ensure `SoftwareInstaller` results don't include the token in error messages or logs. |
| Configuration file corruption | Medium | Use atomic file writes or existing CLI tools (`composer`, `npm`) for updates. |

## Rollback Plan
- Manually edit `~/.composer/auth.json` and `~/.npmrc` to remove the Gitlab token lines.

## Dependencies
- `github.com/charmbracelet/bubbles` (optional, for better text input).

## Success Criteria
- [ ] `~/.composer/auth.json` has the correct `gitlab-token` entry.
- [ ] `~/.npmrc` has the correct `_authToken` entry for Gitlab.
- [ ] The TUI workflow feels seamless when entering the token.
