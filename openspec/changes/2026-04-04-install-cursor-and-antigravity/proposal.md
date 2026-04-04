# Proposal: Add Cursor and Google Antigravity IDEs

## Intent
Expand the set of supported IDEs in `so-install` to include **Cursor** and **Google Antigravity**. These are high-productivity tools that complement the existing VS Code support.

## Scope

### In Scope
- Define `Cursor` and `Antigravity` in `internal/core/domain/software.go`.
- Implement `CursorInstaller` in `internal/infrastructure/cursor/`.
- Implement `AntigravityInstaller` in `internal/infrastructure/antigravity/`.
- Update `cmd/so-install/main.go` to register both.
- Add unit tests for both installers.

### Out of Scope
- Support for non-Debian distributions for these specific tools (currently focusing on `apt`).
- Configuration or plugin installation for these IDEs.

## Approach
- **Cursor**: Download the `.deb` package from `https://downloader.cursor.sh/linux/debian/amd64` to `/tmp/cursor.deb` and install using `apt install`.
- **Antigravity**: Add the Google repo key and source list, then `apt update` and `apt install antigravity`. This installs both the IDE and the `agy` CLI.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Added constants and updated `GetSteps`. |
| `internal/infrastructure/cursor/` | New | New installer implementation. |
| `internal/infrastructure/antigravity/` | New | New installer implementation. |
| `cmd/so-install/main.go` | Modified | Registration of new installers. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Network Failure | Med | Standard error handling for `wget`/`curl`. |
| Apt Lock | Med | Already handled in `InstallSoftwareUseCase` with retries. |
| Repo key change | Low | Documentation of the source for easy updates. |

## Rollback Plan
- Manual: `sudo apt remove cursor antigravity` and remove `/etc/apt/sources.list.d/antigravity.list`.
- Automated: Reverting the code changes will prevent future installs, but won't uninstall existing ones (consistent with current project behavior).

## Dependencies
- `wget`, `curl`, `gnupg` (provided by `BaseDeps`).

## Success Criteria
- [ ] `so-install` lists Cursor and Google Antigravity in the "ides" section.
- [ ] Selecting and installing them works without errors.
- [ ] `cursor` and `agy` commands are available in the shell after install.
