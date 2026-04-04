## Exploration: Install VS Code

### Current State
The system has several installation steps (system-prep, browsers, docker, etc.) but no category for IDEs or development environments beyond CLI tools. Software is registered in `internal/core/domain/software.go` and implemented in `internal/infrastructure/`.

### Affected Areas
- `internal/core/domain/software.go` — Add `VsCode` constant, update `GetSteps` and `AllSoftware`.
- `internal/core/domain/software_test.go` — Add test case for VS Code display name.
- `internal/infrastructure/vscode/vscode.go` — New installer implementation.
- `internal/infrastructure/vscode/vscode_test.go` — New tests for the installer.
- `cmd/so-install/main.go` — Register the new installer.
- `PRD.md` — (Optional) Update to include IDE module.

### Approaches
1. **Official .deb Package** — Download the `.deb` from Microsoft and install it via `apt`.
   - Pros: Simple, consistent with `ChromeInstaller`, handles repo registration automatically on install.
   - Cons: Requires `wget` (already in `BaseDeps`).
   - Effort: Low

2. **Flatpak** — Install `com.visualstudio.code` from Flathub.
   - Pros: Isolated, easy to uninstall.
   - Cons: `code` CLI tool might not be available in the PATH without extra config, harder to integrate with system tools (SDKs, compilers).
   - Effort: Low (using existing `FlatpakAppInstaller`)

### Recommendation
Use the **Official .deb Package** approach. Developers usually prefer the native package for better integration with the terminal and other system tools. It follows the pattern established for browsers like Google Chrome.

### Risks
- Download might fail if the URL changes.
- Redirects: `wget` must handle redirects (default behavior).
- Filename: Use `wget -O /tmp/vscode.deb` to ensure a consistent filename.

### Ready for Proposal
Yes — I have a clear plan for implementation following the project's architecture and patterns.
