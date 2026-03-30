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
