## Exploration: Screen Lock Auto-configuration

### Current State
The project manages software installations and basic configurations for Go/Debian-based systems. It currently lacks functionality to configure desktop-specific settings such as screen lock timeouts.

### Affected Areas
- `internal/core/domain/software.go` — Add `ScreenLockConfig` to `SoftwareID`.
- `internal/core/domain/os_detector.go` — Add `DesktopEnvironment` detection (already exists).
- `internal/infrastructure/desktop/screen_lock.go` — (New) Implementation of the screen lock config for GNOME and KDE.
- `cmd/so-install/main.go` — Register the new installer in the use case.

### Approaches
1. **Infrastructure implementation with DE detection**
   - Use `gsettings` for GNOME:
     - `idle-delay`: 900 (15 minutes)
     - `lock-delay`: 15 (15 seconds)
     - `lock-enabled`: true
   - Use `kwriteconfig5` for KDE:
     - `Timeout`: 900 (15 minutes)
     - `LockGrace`: 15 (15 seconds)
     - `Lock`: true
   - Execution strategy: Since the app might run as root (via `sudo`), we must ensure these commands run as the actual user to modify their session settings. We can use `sudo -u $SUDO_USER` if running as root.

   - Pros: Consistent with current architecture, easy to test with mocks.
   - Cons: Desktop environments might have variations (e.g., Plasma 5 vs 6).
   - Effort: Low

### Recommendation
Use Approach 1. It integrates well with the existing `SoftwareInstaller` pattern and allows for easy expansion if other DEs are supported in the future.

### Risks
- **DBus Session**: `gsettings` and some KDE commands require an active DBus session. Running them from a `sudo` context might require setting `DBUS_SESSION_BUS_ADDRESS`.
- **KDE Versions**: Differences between `kwriteconfig5` and `kwriteconfig6`. We should check which one is available.

### Ready for Proposal
Yes. The feature is well-defined and feasible within the current architecture.
