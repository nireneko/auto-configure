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
