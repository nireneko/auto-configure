# Delta for TUI: Gitlab Configuration

## ADDED Requirements

### Requirement: Token Input Screen
The TUI MUST provide an interactive screen to capture the Gitlab token if `gitlab-token-config` is selected.

#### Scenario: Display token input
- GIVEN the user has selected `gitlab-token-config`
- WHEN the user confirms the selection
- THEN the TUI MUST transition to a token input state before starting the installation.

### Requirement: Masked Token Input
The TUI SHOULD mask the token characters as they are typed to ensure security.

#### Scenario: Hide token characters
- GIVEN the user is typing the Gitlab token
- WHEN each character is entered
- THEN the TUI MUST display a placeholder character (e.g., `*`) instead of the actual character.
