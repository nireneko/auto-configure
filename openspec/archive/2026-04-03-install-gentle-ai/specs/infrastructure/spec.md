# Delta for Infrastructure: Gentle-AI

## ADDED Requirements

### Requirement: Gentle-AI Installer
The system MUST provide an installer for the `gentle-ai` ID that downloads and runs the official installation script.

#### Scenario: Install Gentle-AI
- GIVEN a Debian-based system
- WHEN the `GentleAIInstaller` is executed
- THEN it MUST run `curl -fsSL https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install.sh | bash`
- AND the execution MUST be successful (exit code 0).

### Requirement: User Context Execution for Gentle-AI
The installer MUST execute the installation script as the actual user (the one who invoked sudo).

#### Scenario: Execute as actual user
- GIVEN the application is running as root via sudo
- WHEN `Install()` is called
- THEN the command MUST be executed via `sudo -u $SUDO_USER bash -c ...`
- AND the installer MUST use `domain.GetActualUser()` to find the user.

### Requirement: Gentle-AI Verification
The installer MUST verify the installation by checking if the `gentle-ai` binary is available in the user's path or by checking its version.

#### Scenario: Verify Gentle-AI installation
- GIVEN `Gentle-AI` is installed
- WHEN `IsInstalled()` is called
- THEN it MUST run `gentle-ai --version`
- AND it MUST return `true` if the command succeeds.
