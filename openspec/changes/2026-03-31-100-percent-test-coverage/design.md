# Design: 100% Test Coverage

## Technical Approach

We will systematically add test files or expand existing ones to cover the remaining untested statements, aiming for exactly 100.0% coverage project-wide. Where necessary, code will be slightly refactored to improve testability (e.g., extracting the core logic of `main()` to a `run()` function).

## Architecture Decisions

### Decision: Main execution testability

**Choice**: Refactor `cmd/so-install/main.go` to use a `run(args []string) error` pattern.
**Alternatives considered**: Executing the binary as a subprocess in tests.
**Rationale**: Subprocess testing requires compiling the binary and makes coverage tracking harder. Refactoring to a testable function allows pure unit testing.

### Decision: TUI State Testing

**Choice**: Direct `Update()` and `View()` invocation testing.
**Alternatives considered**: Using `teatest` for full integration testing.
**Rationale**: `teatest` adds external dependency complexity and can be flaky with timing. Testing the state machine by manually passing `tea.KeyMsg` or custom messages into `Update()` and asserting on the resulting `Model` state is deterministic and robust.

### Decision: Testing os/user under sudo

**Choice**: Abstract or mock the environment variables/lookup. Since `GetActualUser` heavily relies on `os/user.Current()` and `os.Getenv`, we can manipulate environment variables in the test using `t.Setenv()`.
**Alternatives considered**: Injecting a user lookup interface.
**Rationale**: Environment variables (`SUDO_USER`) are natively supported in Go tests via `t.Setenv()`, keeping the domain logic simple.

## Data Flow

No architectural data flow changes. This purely adds test vectors acting on existing functions.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `cmd/so-install/main.go` | Modify | Extract `main` body into `run()` for testability. |
| `cmd/so-install/main_test.go` | Create | Test `run()` execution flows. |
| `internal/core/domain/errors_test.go` | Create | Cover `WrapInstallError` and error string formats. |
| `internal/core/domain/user_test.go` | Create | Cover `GetActualUser` and `GetActualHome` with mock env vars. |
| `internal/infrastructure/osrelease/detector_test.go` | Modify | Expand scenarios for `detectDesktopEnvironment` and `isProcessRunning`. |
| `internal/infrastructure/gitlab/configurator_test.go` | Modify | Cover missing Composer and NPM config edge cases. |
| `internal/presentation/tui/model_test.go` | Modify | Ensure all view states and step transitions are tested. |
| `internal/infrastructure/*/*_test.go` | Modify | Add simple tests verifying `ID()` output on installers. |
| `pkg/mocks/software_installer_test.go` | Create | Simple test to instantiate mock and verify basic coverage of its ID function. |

## Interfaces / Contracts

No new interfaces introduced. Existing `SoftwareInstaller`, `Executor`, and `OSDetector` mocks will be utilized heavily.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | All layers | `go test -cover ./...`, using `t.Setenv()` for env and mock executors for shell commands. |

## Migration / Rollout

No migration required.

## Open Questions

- [ ] Will CI environments properly report the 100% coverage without issue? (Assume yes)