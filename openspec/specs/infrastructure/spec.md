# Delta for Infrastructure: System Preparation

## ADDED Requirements

### Requirement: System Update Installer
The system MUST provide an installer for the `system-update` ID that updates the system's package list and upgrades existing packages.

#### Scenario: Update and Upgrade
- GIVEN a Debian-based system
- WHEN the `SystemUpdateInstaller` is executed
- THEN it MUST run `apt-get update`
- AND it MUST run `apt-get upgrade -y` with `DEBIAN_FRONTEND=noninteractive`

### Requirement: Base Dependencies Installer
The system MUST provide an installer for the `base-deps` ID that installs essential system tools.

#### Scenario: Install base tools
- GIVEN a Debian-based system
- WHEN the `BaseDepsInstaller` is executed
- THEN it MUST install `git`, `wget`, `curl`, `ca-certificates`, `gnupg`, and `lsb-release`
- AND it MUST use the `-y` flag for non-interactive installation

### Requirement: Gitlab Token Configurator
The system MUST provide a configurator for the `gitlab-token-config` ID that updates global Composer and NPM settings with a Gitlab personal access token.

### Requirement: Composer Global Auth
The Gitlab configurator MUST update the user's `~/.composer/auth.json` file with the Gitlab token for `gitlab.com`.

#### Scenario: Update Composer auth.json
- GIVEN a valid Gitlab token
- WHEN the Gitlab configurator is executed
- THEN it MUST ensure `~/.composer` directory exists
- AND it MUST add or update the `gitlab-token` entry for `gitlab.com` in `auth.json`

### Requirement: NPM Global Auth
The Gitlab configurator MUST update the user's `~/.npmrc` file with the Gitlab token for `gitlab.com`.

#### Scenario: Update npmrc
- GIVEN a valid Gitlab token
- WHEN the Gitlab configurator is executed
- THEN it MUST append or update `//gitlab.com/api/v4/packages/npm/:_authToken=TOKEN` in `~/.npmrc`

### Requirement: GNOME Screen Lock Configuration

The system MUST configure GNOME screen lock settings when the detected desktop environment is GNOME.

#### Scenario: Apply GNOME Screen Lock
- GIVEN the detected desktop environment is GNOME
- WHEN the screen lock configuration is applied
- THEN `gsettings set org.gnome.desktop.session idle-delay 900` MUST be executed
- AND `gsettings set org.gnome.desktop.screensaver lock-delay 15` MUST be executed
- AND `gsettings set org.gnome.desktop.screensaver lock-enabled true` MUST be executed.

### Requirement: KDE Screen Lock Configuration

The system MUST configure KDE screen lock settings when the detected desktop environment is KDE.

#### Scenario: Apply KDE Screen Lock
- GIVEN the detected desktop environment is KDE
- WHEN the screen lock configuration is applied
- THEN `kwriteconfig5` or `kwriteconfig6` MUST be used to set `Timeout` to 900 and `LockGrace` to 15 in `kscreenlockerrc`.
- AND a DBus notification MUST be sent to reload the configuration.

### Requirement: User Context Execution

The system MUST execute desktop configuration commands in the context of the actual user.

#### Scenario: Execute as actual user
- GIVEN the application is running as root via sudo
- WHEN a desktop configuration command is executed
- THEN the command MUST be executed as the original user (found via `SUDO_USER`).
