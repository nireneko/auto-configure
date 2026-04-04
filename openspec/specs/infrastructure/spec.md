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

### Requirement: Test Edge Cases for Infrastructure Components
The system MUST include unit tests that cover edge cases y fallos de ejecución en `openvpn`, `nvm` y `homebrew`.

#### Scenario: Fallo de ejecución del instalador de OpenVPN
- GIVEN que el ejecutor del sistema está mockeado
- WHEN el instalador de OpenVPN intenta descargar el script pero la red falla
- THEN el instalador debe retornar un error explícito de red

#### Scenario: Fallo de ejecución del instalador de NVM
- GIVEN que el ejecutor del sistema está mockeado
- WHEN el instalador de NVM intenta descargar el script de instalación
- THEN si el comando `wget` falla, el instalador debe retornar el error correspondiente

#### Scenario: Fallo en la verificación de Homebrew
- GIVEN que Homebrew no está instalado
- WHEN el instalador verifica si el comando `brew` existe
- THEN si el comando no existe y la instalación falla, debe retornar error

### Requirement: GitlabTokenConfigurator — Composer File Ownership

The configurator MUST chown the `~/.composer/` directory and `~/.composer/auth.json` to the real user after creating/writing them.

#### Scenario: Install under sudo — composer files chowned to real user

- GIVEN `SUDO_UID=1000`, `SUDO_GID=1000`, and a `chownFn` spy injected
- WHEN `Install()` is called with a valid token
- THEN `chownFn` is called with `("~/.composer", 1000, 1000)`
- AND `chownFn` is called with `("~/.composer/auth.json", 1000, 1000)`

#### Scenario: Install without sudo — composer files chowned to process owner

- GIVEN `SUDO_UID` and `SUDO_GID` are not set, and a `chownFn` spy injected
- WHEN `Install()` is called with a valid token
- THEN `chownFn` is called with the current process UID/GID for both composer paths

### Requirement: GitlabTokenConfigurator — NPM File Ownership

The configurator MUST chown `~/.npmrc` to the real user after writing it.

#### Scenario: Install under sudo — .npmrc chowned to real user

- GIVEN `SUDO_UID=1000`, `SUDO_GID=1000`, and a `chownFn` spy injected
- WHEN `Install()` is called with a valid token
- THEN `chownFn` is called with `("~/.npmrc", 1000, 1000)`

#### Scenario: Install without sudo — .npmrc chowned to process owner

- GIVEN `SUDO_UID` and `SUDO_GID` are not set, and a `chownFn` spy injected
- WHEN `Install()` is called with a valid token
- THEN `chownFn` is called with the current process UID/GID for `.npmrc`

### Requirement: GitlabTokenConfigurator — chownFn Injection

The configurator MUST accept an injectable `chownFn func(string, int, int) error` (defaulting to `os.Chown`) and injectable `uidFn`/`gidFn` functions. A `SetChownFn` setter MUST be provided for test injection.

#### Scenario: Default behavior uses os.Chown

- GIVEN a configurator created with `NewGitlabTokenConfigurator`
- WHEN `Install()` succeeds
- THEN the real `os.Chown` is invoked (no custom fn set)

#### Scenario: chownFn failure is propagated

- GIVEN a `chownFn` spy that returns an error
- WHEN `Install()` is called
- THEN `Install()` returns that error
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
