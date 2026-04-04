# Proposal: Install Gentle-AI after AI tools

## Intent
Add Gentle-AI to the installation process. Gentle-AI supercharges AI agents with memory, skills, and a senior architect persona. It should be available for installation immediately after the base AI CLI tools (Gemini, Claude, Codex, etc.).

## Scope

### In Scope
- Define `GentleAI` SoftwareID in domain.
- Implement `GentleAIInstaller` using the official curl-based script.
- Create a dedicated "gentle-ai" installation phase in the TUI sequence.
- Add automated tests for the new installer.
- Update `main.go` to register the installer.

### Out of Scope
- Automatic configuration of agents using `gentle-ai install` (manual step for the user).
- Support for Windows/macOS (project is currently Linux-focused).

## Approach
1.  **Domain Update**: Add `GentleAI` constant to `SoftwareID`.
2.  **Infrastructure**: Create `internal/infrastructure/gentleai` with `GentleAIInstaller` that runs the official install script.
3.  **TUI/Steps**: Insert a new `InstallStep` for Gentle-AI between `ai-cli` and `flatpak`.
4.  **Registration**: Initialize the installer in `cmd/so-install/main.go`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Add `GentleAI` ID and new step in `GetSteps`. |
| `internal/infrastructure/gentleai` | New | Implementation and tests for Gentle-AI installer. |
| `cmd/so-install/main.go` | Modified | Registration of the new installer. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Network failure during curl | Medium | Proper error wrapping and reporting in TUI. |
| Script requires user interaction | Low | Use `-fsSL` and assume non-interactive mode. |
| Binary not in PATH after install | Low | Usually handled by the install script (adds to .bashrc). |

## Rollback Plan
Remove the newly created files and revert changes in `software.go` and `main.go`. The tool itself can be uninstalled manually if needed (reverting `.bashrc` changes).

## Dependencies
- `curl` and `bash` (already part of `base-deps`).

## Success Criteria
- [ ] `Gentle-AI` appears in the TUI software selection list.
- [ ] Installation successfully downloads and installs the `gentle-ai` binary.
- [ ] `gentle-ai --version` returns success after installation.
- [ ] Tests for the installer pass with 100% coverage.
