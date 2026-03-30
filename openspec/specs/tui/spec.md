# Delta for TUI: System Preparation

## MODIFIED Requirements

### Requirement: Mandatory Software Selection
The TUI MUST ensure that `system-update` and `base-deps` are always included in the installation sequence, even if not selectable in the UI.

#### Scenario: Prepending mandatory steps
- GIVEN the user has selected one or more software items from the list
- WHEN the user confirms the selection to start the installation
- THEN the system MUST prepend `SystemUpdate` and `BaseDeps` to the selected software list
- AND these steps MUST be visible in the progress view
