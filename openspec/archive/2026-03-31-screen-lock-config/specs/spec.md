# Delta for Domain

## ADDED Requirements

### Requirement: Screen Lock Configuration

The system SHALL provide a way to configure the screen lock and idle timeout for the user's desktop environment.

#### Scenario: Screen Lock Software ID
- GIVEN a list of all supported software
- WHEN the user views the available software
- THEN "Screen Lock Configuration" MUST be included in the list.

# Delta for Infrastructure

## ADDED Requirements

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
