# Proposal: Install Bitwarden via Flatpak

## Intent
Add Bitwarden to the list of installable software using the Flatpak infrastructure. This fulfills the user's need for a password manager and leverages the recently added Flatpak support.

## Scope

### In Scope
- Add `Bitwarden` to `domain.SoftwareID` and `domain.GetSteps()`.
- Implement a generic `FlatpakAppInstaller` in `internal/infrastructure/flatpak/app_installer.go`.
- Wire `Bitwarden` in `main.go` using the generic installer.
- Add unit tests for the new installer.

### Out of Scope
- Installing Bitwarden via other methods (snap, deb, etc.).
- Managing Bitwarden configuration or vault.

## Approach
We will create a generic `FlatpakAppInstaller` that implements `domain.SoftwareInstaller`. It will take a `domain.SoftwareID` and a Flatpak application ID (e.g., `com.bitwarden.desktop`). It will use `flatpak install -y flathub <id>` for installation and `flatpak info <id>` to check if it's installed.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Add `Bitwarden` ID and update steps. |
| `internal/infrastructure/flatpak/app_installer.go` | New | Generic Flatpak app installer. |
| `cmd/so-install/main.go` | Modified | Instantiate and wire Bitwarden installer. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Flatpak not installed | Low | Ensure `Flatpak` step runs before `Bitwarden`. |
| Flathub repo missing | Low | `FlatpakInstaller` adds it; verify in `AppInstaller`. |

## Rollback Plan
Remove the added code and revert `main.go` and `domain/software.go`.

## Dependencies
- `flatpak` package and `flathub` repository must be installed/configured.

## Success Criteria
- [ ] Bitwarden appears in the TUI list.
- [ ] Selecting Bitwarden installs the flatpak correctly.
- [ ] If already installed, it's detected correctly.
