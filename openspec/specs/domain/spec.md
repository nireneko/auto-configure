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
