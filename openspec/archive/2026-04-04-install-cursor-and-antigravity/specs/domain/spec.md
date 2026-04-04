# Delta for Domain: IDEs (Cursor & Antigravity)

## ADDED Requirements

### Requirement: Cursor Software ID
The system MUST include `cursor` as a valid `SoftwareID` value.

#### Scenario: Cursor ID registration
- GIVEN the `domain.SoftwareID` type
- WHEN checking for the new ID
- THEN `Cursor` MUST be `"cursor"`

### Requirement: Antigravity Software ID
The system MUST include `antigravity` as a valid `SoftwareID` value.

#### Scenario: Antigravity ID registration
- GIVEN the `domain.SoftwareID` type
- WHEN checking for the new ID
- THEN `Antigravity` MUST be `"antigravity"`

### Requirement: IDEs Step Update
The system MUST include `Cursor` and `Antigravity` in the `ides` step of the installation sequence.

#### Scenario: Add Cursor and Antigravity to IDEs step
- GIVEN the `domain.GetSteps()` function
- WHEN retrieving the installation steps
- THEN the step with ID `"ides"` MUST contain `VsCode`, `Cursor`, and `Antigravity`

### Requirement: IDEs Display Names
The system MUST provide human-readable display names for Cursor and Antigravity.

#### Scenario: Cursor display name
- GIVEN the `domain.Cursor` software ID
- WHEN `DisplayName()` is called
- THEN it MUST return `"Cursor IDE"`

#### Scenario: Antigravity display name
- GIVEN the `domain.Antigravity` software ID
- WHEN `DisplayName()` is called
- THEN it MUST return `"Google Antigravity"`
