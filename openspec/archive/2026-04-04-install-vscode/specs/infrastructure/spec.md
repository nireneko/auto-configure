# Delta for Infrastructure: IDEs

## ADDED Requirements

### Requirement: VS Code Installer
The system MUST provide an installer for the `vscode` ID that downloads and installs the official `.deb` package.

#### Scenario: Install VS Code
- GIVEN a Debian-based system
- WHEN the `VsCodeInstaller` is executed
- THEN it MUST download the latest `.deb` from `https://go.microsoft.com/fwlink/?LinkID=760868`
- AND it MUST save it to `/tmp/vscode.deb`
- AND it MUST run `apt install -y /tmp/vscode.deb`

### Requirement: VS Code Verification
The installer MUST verify the installation by checking if the `code` binary exists in the system path.

#### Scenario: Verify VS Code installation
- GIVEN `code` is installed
- WHEN `IsInstalled()` is called
- THEN it MUST run `which code`
- AND it MUST return `true` if the command succeeds.

### Requirement: VS Code Architecture
The installer SHOULD target the `amd64` architecture, as it is the standard for the target systems of this project.

#### Scenario: AMD64 Architecture Link
- GIVEN the installer download logic
- WHEN constructing the wget command
- THEN it MUST use the link `https://go.microsoft.com/fwlink/?LinkID=760868` (which redirects to the 64-bit .deb)
