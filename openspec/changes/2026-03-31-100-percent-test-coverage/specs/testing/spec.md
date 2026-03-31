# Testing Specification

## Purpose

This specification defines the test coverage requirements across all domains of the 1x-so-install application to ensure complete coverage, including edge cases.

## Requirements

### Requirement: 100% Test Coverage

The system MUST have 100% statement coverage as reported by `go test -cover`. 

#### Scenario: Full suite execution
- GIVEN the test suite is run with `make test` or `go test -cover ./...`
- WHEN the tests finish executing
- THEN the coverage tool reports 100.0% coverage for all packages.

### Requirement: TUI Event Handling Testing

The TUI presentation layer MUST be tested for all key messages, state transitions, and view renderings.

#### Scenario: User navigation
- GIVEN the application is in `stateWelcome`
- WHEN the user presses `Enter`
- THEN the model updates to `stateSoftwareSelect`.

### Requirement: OS and Home Directory Retrieval Edge Cases

The system MUST handle environment variables like `SUDO_USER` when resolving paths or privileges.

#### Scenario: Executed under sudo
- GIVEN `SUDO_USER` is set
- WHEN `GetActualUser` is called
- THEN it returns the user information for `SUDO_USER`, not `root`.

#### Scenario: Installers ID checking
- GIVEN an initialized infrastructure installer
- WHEN `ID()` is called
- THEN it MUST return its corresponding `domain.SoftwareID` without panicking.