# Delta for Domain: Gitlab Configuration

## ADDED Requirements

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
