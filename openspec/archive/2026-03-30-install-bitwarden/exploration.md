## Exploration: Bitwarden via Flatpak

### Current State
The system has a `FlatpakInstaller` that installs the `flatpak` package and adds the `flathub` repository. However, there is no mechanism to install specific applications from Flathub yet.

### Affected Areas
- `internal/core/domain/software.go` — Add `Bitwarden` SoftwareID, display name, and update steps.
- `internal/infrastructure/flatpak/app_installer.go` — (New file) Generic installer for Flatpak applications.
- `cmd/so-install/main.go` — Wire the new `Bitwarden` installer.

### Approaches
1. **Generic FlatpakAppInstaller (Recommended)** — Similar to `NpmInstaller`, create a struct that takes the Flatpak ID (e.g., `com.bitwarden.desktop`).
   - Pros: Reusable for future apps (Slack, Discord, etc.), consistent with the NPM pattern.
   - Cons: None.
   - Effort: Low.

2. **Specific BitwardenInstaller** — Hardcoded installer for Bitwarden.
   - Pros: Simple.
   - Cons: Not reusable, violates DRY if we add more Flatpak apps.
   - Effort: Low.

### Recommendation
I recommend **Option 1**. It allows us to scale the list of available software easily. We will create `FlatpakAppInstaller` which uses `flatpak install -y flathub <app-id>`.

### Risks
- **Dependency**: Bitwarden requires Flatpak and Flathub to be present. We must ensure the `Flatpak` step runs before the `Bitwarden` step.
- **Flathub availability**: Requires internet connection and Flathub to be reachable (already a risk for other installers).

### Ready for Proposal
Yes.
