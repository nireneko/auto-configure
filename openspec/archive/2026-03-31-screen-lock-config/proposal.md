# Proposal: Screen Lock Auto-configuration

## Intent
Add automatic configuration for screen lock and idle timeout to enhance system security and energy efficiency out-of-the-box. This addresses the need for a standardized "secure" default for GNOME and KDE environments.

## Scope

### In Scope
- Detection of GNOME and KDE desktop environments.
- Setting idle timeout to 15 minutes (900s).
- Setting screen lock delay to 15 seconds.
- Enabling screen lock.
- Ensuring commands run as the actual user (handling `sudo` context).

### Out of Scope
- Support for other desktop environments (XFCE, Sway, etc.) at this stage.
- Customizing the timeouts via CLI arguments (fixed for now as requested).

## Approach
Implement a new `ScreenLockInstaller` in `internal/infrastructure/desktop` that satisfies the `domain.SoftwareInstaller` interface. It will use `gsettings` for GNOME and `kwriteconfig5`/`kwriteconfig6` for KDE. To handle `sudo` correctly, it will use `domain.GetActualUser()` and wrap commands with `sudo -u $USER` when necessary.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modified | Add `ScreenLockConfig` SoftwareID. |
| `internal/infrastructure/desktop/` | New | New package for desktop environment configurations. |
| `cmd/so-install/main.go` | Modified | Register the new installer. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| DBus session missing | Low | Ensure commands run in the user context. |
| KDE version mismatch | Medium | Check for `kwriteconfig6` first, fallback to `kwriteconfig5`. |

## Rollback Plan
Manual reset via DE settings or running inverse commands. No automatic backup of previous values.

## Success Criteria
- [ ] GNOME: `idle-delay` is 900, `lock-delay` is 15.
- [ ] KDE: `Timeout` is 900, `LockGrace` is 15.
- [ ] Installation step completes without error in both environments.
