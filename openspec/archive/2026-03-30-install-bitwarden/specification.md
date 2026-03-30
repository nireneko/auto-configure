# Software Installation: Bitwarden via Flatpak

## Purpose
Define the requirements and scenarios for installing Bitwarden on the system using the Flatpak package manager.

## Requirements

### Requirement: Bitwarden Software ID
The system MUST include `bitwarden` as a valid `SoftwareID` in the domain.

#### Scenario: Bitwarden ID exists
- GIVEN the domain software list
- WHEN checking for the Bitwarden ID
- THEN the constant `Bitwarden` MUST be present and equal to `"bitwarden"`

### Requirement: Flatpak App Installer
The system MUST provide a generic installer for Flatpak applications from the Flathub repository.

#### Scenario: Install Bitwarden via Flatpak
- GIVEN Flatpak is installed and Flathub is configured
- WHEN the user selects Bitwarden for installation
- THEN the system MUST execute `flatpak install -y flathub com.bitwarden.desktop`

#### Scenario: Detect already installed Bitwarden
- GIVEN Bitwarden is already installed via Flatpak
- WHEN checking the installation status
- THEN the system MUST report it as already installed using `flatpak info com.bitwarden.desktop`

### Requirement: Installation Step Sequencing
Bitwarden MUST be positioned in the installation sequence AFTER the Flatpak runtime installation.

#### Scenario: Sequencing check
- GIVEN the list of installation steps
- WHEN examining the order
- THEN the `bitwarden` software MUST appear in a step that follows the `flatpak` runtime step

### Requirement: Error Handling
The system MUST handle installation failures gracefully, providing the user with relevant error information.

#### Scenario: Flatpak installation fails
- GIVEN the `flatpak install` command returns an error
- WHEN the installation process is running
- THEN the system MUST wrap the error with context and return it as an `InstallResult`
