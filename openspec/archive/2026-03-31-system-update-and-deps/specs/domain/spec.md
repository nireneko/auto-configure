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
