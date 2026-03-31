# Delta for TUI

## MODIFIED Requirements

### Requirement: TUI Testing Compatibility
The TUI tests MUST be able to compile and correctly assert messages and state against the `Update` method of Bubbletea.
(Previously: The tests were asserting single return values from `Update` which now returns two values: `tea.Model, tea.Cmd`).

#### Scenario: Update Method Usage in Tests
- GIVEN a mocked tea.Model
- WHEN a message is passed to `Update()`
- THEN the tests must assign the return value to two variables (`m, cmd`) and assert against `m` accordingly without compilation errors.