# Testing Specification: 90% Coverage & Edge Cases

## Purpose
Ensure the robustness of the `so-install` tool by reaching 90% test coverage and validating critical edge cases.

## Requirements

### Requirement: 90% Minimum Coverage
Every package MUST have at least 90% statement coverage.

#### Scenario: Full suite execution
- GIVEN the test suite is run
- WHEN the tests finish
- THEN every package reports >= 90% coverage.

### Requirement: Testable Entry Point
The main entry point MUST be testable without running the actual binary.

#### Scenario: Run with --help
- GIVEN the application logic is called with `[]string{"--help"}`
- WHEN the `Run` function executes
- THEN it returns a success exit code and prints help (if applicable).

### Requirement: Desktop Environment Detection
The system MUST correctly identify KDE and Gnome environments.

#### Scenario: KDE Detection
- GIVEN a system with `KDE_FULL_SESSION` set or `plasmashell` running
- WHEN `detectDesktopEnvironment` is called
- THEN it returns `KDE`.

### Requirement: Nvidia Installation Logic
All Nvidia installation branches (Free vs Proprietary) MUST be exercised.

#### Scenario: Proprietary Nvidia Install
- GIVEN an Nvidia installer on a non-Debian system
- WHEN `Install` is called with proprietary preference
- THEN it executes the `ubuntu-drivers` or similar proprietary logic.

### Requirement: TUI State Transitions
All TUI states MUST be reachable and renderable in tests.

#### Scenario: Nvidia Config to Next State
- GIVEN the TUI is in `stateNvidiaConfig`
- WHEN the user selects an option
- THEN it transitions to the next logical state (`stateSoftwareSelect` or `stateTokenInput`).
