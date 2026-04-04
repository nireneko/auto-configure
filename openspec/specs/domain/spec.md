# Delta for Domain: System Preparation

## ADDED Requirements

### Requirement: System Preparation Software IDs
The system MUST include `system-update` and `base-deps` as valid `SoftwareID` values.

#### Scenario: Software IDs registration
- GIVEN the `domain.SoftwareID` type
- WHEN checking for the new IDs
- THEN `SystemUpdate` MUST be `"system-update"`
- AND `BaseDeps` MUST be `"base-deps"`

### Requirement: System Preparation Step
The system MUST include a "system-prep" step at the beginning of the installation sequence.

#### Scenario: First step is system-prep
- GIVEN the `domain.GetSteps()` function
- WHEN retrieving the installation steps
- THEN the first step MUST have ID `"system-prep"`
- AND it MUST contain `SystemUpdate` and `BaseDeps`
- AND it MUST be marked as `Critical: true`

### Requirement: Gitlab Config Software ID
The system MUST include `gitlab-token-config` as a valid `SoftwareID` value.

#### Scenario: Gitlab Config ID registration
- GIVEN the `domain.SoftwareID` type
- WHEN checking for the new ID
- THEN `GitlabTokenConfig` MUST be `"gitlab-token-config"`

### Requirement: Apps Step Update
The system MUST include `GitlabTokenConfig` in the `apps` step of the installation sequence.

#### Scenario: Add Gitlab to Apps step
- GIVEN the `domain.GetSteps()` function
- WHEN retrieving the installation steps
- THEN the step with ID `"apps"` MUST contain `GitlabTokenConfig`

### Requirement: Screen Lock Configuration

The system SHALL provide a way to configure the screen lock and idle timeout for the user's desktop environment.

#### Scenario: Screen Lock Software ID
- GIVEN a list of all supported software
- WHEN the user views the available software
- THEN "Screen Lock Configuration" MUST be included in the list.

### Requirement: GetActualUID

The system MUST expose `GetActualUID() int` in `domain/user.go`. It MUST return the integer value of `SUDO_UID` when that env var is set and parseable; otherwise it MUST return `os.Getuid()`.

#### Scenario: SUDO_UID is set to a valid integer

- GIVEN `SUDO_UID` env var is set to `"1000"`
- WHEN `GetActualUID()` is called
- THEN it returns `1000`

#### Scenario: SUDO_UID is not set

- GIVEN `SUDO_UID` env var is empty or absent
- WHEN `GetActualUID()` is called
- THEN it returns the current process UID (`os.Getuid()`)

#### Scenario: SUDO_UID is set to a non-integer value

- GIVEN `SUDO_UID` env var is set to `"abc"`
- WHEN `GetActualUID()` is called
- THEN it returns the current process UID (parse fallback)

### Requirement: GetActualGID

The system MUST expose `GetActualGID() int` in `domain/user.go`. It MUST return the integer value of `SUDO_GID` when set and parseable; otherwise it MUST return `os.Getgid()`.

#### Scenario: SUDO_GID is set to a valid integer

- GIVEN `SUDO_GID` env var is set to `"1000"`
- WHEN `GetActualGID()` is called
- THEN it returns `1000`

#### Scenario: SUDO_GID is not set

- GIVEN `SUDO_GID` env var is empty or absent
- WHEN `GetActualGID()` is called
- THEN it returns the current process GID (`os.Getgid()`)

#### Scenario: SUDO_GID is set to a non-integer value

- GIVEN `SUDO_GID` env var is set to `"abc"`
- WHEN `GetActualGID()` is called
- THEN it returns the current process GID (parse fallback)
# Delta for Domain: Gentle-AI

## ADDED Requirements

### Requirement: Gentle-AI Software ID
The system MUST include `gentle-ai` as a valid `SoftwareID` value.

#### Scenario: Gentle-AI ID registration
- GIVEN the `domain.SoftwareID` type
- WHEN checking for the new ID
- THEN `GentleAI` MUST be `"gentle-ai"`

### Requirement: Gentle-AI Install Step
The system MUST include a "gentle-ai" step in the installation sequence, positioned immediately after the "ai-cli" step.

#### Scenario: Gentle-AI step position
- GIVEN the `domain.GetSteps()` function
- WHEN retrieving the installation steps
- THEN the step with ID `"gentle-ai"` MUST be found
- AND its index MUST be `ai-cli_index + 1`
- AND it MUST contain only the `GentleAI` software ID
- AND it MUST be marked as `Critical: false`
# Delta for Domain: IDEs

## ADDED Requirements

### Requirement: VS Code Software ID
The system MUST include `vscode` as a valid `SoftwareID` value.

#### Scenario: VS Code ID registration
- GIVEN the `domain.SoftwareID` type
- WHEN checking for the new ID
- THEN `VsCode` MUST be `"vscode"`

### Requirement: IDEs Install Step
The system MUST include an "ides" step in the installation sequence, positioned after the "gentle-ai" step.

#### Scenario: IDEs step position
- GIVEN the `domain.GetSteps()` function
- WHEN retrieving the installation steps
- THEN a step with ID `"ides"` MUST be found
- AND its index MUST be `gentle-ai_index + 1`
- AND it MUST contain the `VsCode` software ID
- AND it MUST be marked as `Critical: false`

### Requirement: VS Code Display Name
The system MUST provide a human-readable display name for VS Code.

#### Scenario: VS Code display name
- GIVEN the `domain.VsCode` software ID
- WHEN `DisplayName()` is called
- THEN it MUST return `"Visual Studio Code"`
