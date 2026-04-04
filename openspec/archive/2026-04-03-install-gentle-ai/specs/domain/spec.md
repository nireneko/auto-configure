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
