# Proposal: 100% Test Coverage

## Intent

Achieve 100% test coverage across the entire project (`github.com/so-install/*`) to ensure stability, verify edge cases, and improve overall code quality.

## Scope

### In Scope
- Add unit tests for `cmd/so-install/main.go`
- Add unit tests for `internal/core/domain/errors.go` and `user.go`
- Add unit tests for `internal/infrastructure/osrelease/detector.go`
- Add unit tests for `internal/infrastructure/gitlab/configurator.go`
- Add UI tests for `internal/presentation/tui/model.go` using `teatest` or manual updates
- Add tests for all infrastructure installer `ID()` methods and missed branches
- Fix coverage for `pkg/mocks/software_installer.go`

### Out of Scope
- Refactoring architecture or adding new features
- Integration tests or E2E tests

## Approach

Systematically write table-driven tests for uncovered files and methods using `stretchr/testify` and shared mocks from `pkg/mocks`. For the TUI, use Charm's `teatest` or construct `tea.Msg` events directly to simulate user interactions and system events.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `cmd/so-install/main_test.go` | New | Main entry point test |
| `internal/core/domain/*_test.go` | New/Modified | Domain layer edge cases |
| `internal/infrastructure/*/*_test.go` | Modified | Installer ID and error path tests |
| `internal/presentation/tui/*_test.go` | Modified | TUI state transitions and view logic |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| TUI tests are flaky | Low | Avoid timing-based checks; use deterministic `Update` assertions |
| Testing `os.UserHomeDir` | Low | Stub or mock the environment variables/functions used in `domain/user.go` |

## Rollback Plan

Revert the test files using `git checkout` or `git reset` if tests introduce instability in the CI.

## Dependencies

- None (Standard Go tools + testify + teatest if needed)

## Success Criteria

- [ ] `go test -cover ./...` reports 100.0% coverage for all packages.
- [ ] No race conditions reported (`go test -race`).
