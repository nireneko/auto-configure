# Proposal: Install VS Code

## Intent
Add the option to install Visual Studio Code as the first IDE in the orchestrator. Users need a professional development environment beyond CLI tools.

## Scope

### In Scope
- `VsCode` SoftwareID in domain.
- New `ide` step in installation phases.
- `VsCodeInstaller` using the official `.deb` package.
- Unit tests for the installer and domain updates.
- Registration in `main.go`.

### Out of Scope
- Installation of VS Code extensions.
- Configuration of VS Code settings (beyond default install).
- Other IDEs (IntelliJ, Sublime, etc.) — deferred for future tasks.

## Approach
Implement `SoftwareInstaller` for VS Code. The installer will download the stable `.deb` from Microsoft using `wget` and install it via `apt`. This automatically sets up the Microsoft repository for future updates.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Add `VsCode` ID and `ide` step. |
| `internal/infrastructure/vscode/` | New | `VsCodeInstaller` implementation and tests. |
| `cmd/so-install/main.go` | Modified | Register `VsCodeInstaller`. |
| `PRD.md` | Modified | Add IDE section to functional requirements. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Download failure | Low | Use stable `go.microsoft.com` link; check error in TUI. |
| Apt lock | Med | Handled by existing `domain.AptLockError` pattern. |
| Architecture mismatch | Low | Target `amd64` (standard for this project). |

## Rollback Plan
- Remove the new `internal/infrastructure/vscode` directory.
- Revert changes in `internal/core/domain/software.go` and `cmd/so-install/main.go`.
- If already installed, the user can `sudo apt remove code`.

## Dependencies
- `wget` (provided by `BaseDeps`).
- Internet connection for downloading the `.deb`.

## Success Criteria
- [ ] `VsCode` appears in the TUI software selection list.
- [ ] Selecting `VsCode` successfully installs it.
- [ ] `code --version` works after installation.
- [ ] Unit tests for `VsCodeInstaller` pass.
