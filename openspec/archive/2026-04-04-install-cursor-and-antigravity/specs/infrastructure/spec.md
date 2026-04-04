# Delta for Infrastructure: IDEs (Cursor & Antigravity)

## ADDED Requirements

### Requirement: Cursor Installer
The system MUST provide an installer for the `cursor` ID that downloads and installs the official `.deb` package.

#### Scenario: Install Cursor
- GIVEN a Debian-based system
- WHEN the `CursorInstaller` is executed
- THEN it MUST download the latest `.deb` from `https://downloader.cursor.sh/linux/debian/amd64`
- AND it MUST save it to `/tmp/cursor.deb`
- AND it MUST run `apt install -y /tmp/cursor.deb`

### Requirement: Cursor Verification
The installer MUST verify the installation by checking if the `cursor` binary exists in the system path.

#### Scenario: Verify Cursor installation
- GIVEN `cursor` is installed
- WHEN `IsInstalled()` is called
- THEN it MUST run `which cursor`
- AND it MUST return `true` if the command succeeds.

### Requirement: Antigravity Installer
The system MUST provide an installer for the `antigravity` ID that sets up the official Google repository and installs the package.

#### Scenario: Install Antigravity
- GIVEN a Debian-based system
- WHEN the `AntigravityInstaller` is executed
- THEN it MUST ensure `/etc/apt/keyrings` exists
- AND it MUST download the GPG key from `https://us-central1-apt.pkg.dev/doc/repo-signing-key.gpg` and save it to `/etc/apt/keyrings/antigravity-repo-key.gpg` (de-armored)
- AND it MUST add the repository `deb [signed-by=/etc/apt/keyrings/antigravity-repo-key.gpg] https://us-central1-apt.pkg.dev/projects/antigravity-auto-updater-dev/ antigravity-debian main` to `/etc/apt/sources.list.d/antigravity.list`
- AND it MUST run `apt update`
- AND it MUST run `apt install -y antigravity`

### Requirement: Antigravity Verification
The installer MUST verify the installation by checking if the `antigravity` (IDE) and `agy` (CLI) binaries exist.

#### Scenario: Verify Antigravity installation
- GIVEN `antigravity` is installed
- WHEN `IsInstalled()` is called
- THEN it MUST run `which agy`
- AND it MUST return `true` if the command succeeds.
