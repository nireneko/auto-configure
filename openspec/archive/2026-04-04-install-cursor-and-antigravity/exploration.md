## Exploration: Add Cursor and Google Antigravity IDEs

### Current State
The project supports installing several software packages, including Visual Studio Code (VS Code) as an IDE. The installation pattern for `.deb` packages (VS Code, Chrome) is to download the file to `/tmp/` and then use `apt install`.

### Affected Areas
- `internal/core/domain/software.go` — Add new `SoftwareID` values and register in `GetSteps`.
- `internal/infrastructure/` — Add new packages for `cursor` and `antigravity` installers.
- `cmd/so-install/main.go` — Instantiate and register the new installers in the `installerMap`.

### Approaches
1. **Direct Download (.deb) for Cursor & Repository for Antigravity**
   - **Cursor**: Use the official `.deb` download link: `https://downloader.cursor.sh/linux/debian/amd64`.
   - **Antigravity**: Use the official Google Artifact Registry repository as it's the recommended way to get updates.
   - Pros: Follows official recommendations, integrates well with `apt`.
   - Cons: Slightly different patterns for each, but both use `apt`.
   - Effort: Medium

2. **AppImage for both**
   - Pros: Distribution independent, doesn't require `sudo` for download.
   - Cons: Harder to integrate into the menu and terminal, doesn't follow the project's established `apt` pattern.
   - Effort: High (due to integration)

### Recommendation
**Approach 1** is recommended. It aligns with how VS Code and Chrome are already installed in this project and ensures users get updates via `apt` (especially for Antigravity).

### Risks
- **Antigravity Repo changes**: Since it's in "Public Preview", repo URLs might change.
- **Dependency on `gnupg` and `curl`**: These are required for adding the repository but are already listed in `BaseDeps`.

### Ready for Proposal
Yes. The installation steps are clear and the project structure is ready for these additions.
