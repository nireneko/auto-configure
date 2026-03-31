## Exploration: 100% Test Coverage

### Current State
Current global test coverage is 75.8%. The tests are mostly unit tests utilizing `stretchr/testify` and shared mocks in `pkg/mocks`. The core usecases and many infrastructure installers have good coverage, but several edge cases and specific layers are missing coverage.

### Affected Areas
- `cmd/so-install/main.go` (0%) — Entry point, hard to test without refactoring `main` into a testable `run()` function.
- `internal/core/domain/errors.go` (0%) — Simple error wrapping and error string methods.
- `internal/core/domain/user.go` (0%) — Logic to get actual user/home under `sudo`.
- `internal/infrastructure/osrelease/detector.go` (53.3%) — OS detection logic and desktop environment detection.
- `internal/presentation/tui/model.go` (76.9%) — Bubbletea TUI state machine and key handling.
- `internal/infrastructure/*/` — Various installers are missing coverage for their `ID()` methods or specific error handling paths.
- `pkg/mocks/` — Missing coverage for `ID()` in `software_installer.go`.

### Approaches
1. **Targeted Unit Test Addition** — Add standard `go test` unit tests for all uncovered files.
   - Pros: Standard, follows existing patterns, improves confidence.
   - Cons: TUI testing and `main.go` testing can be tricky.
   - Effort: Medium

### Recommendation
Proceed with adding targeted unit tests, mocking out system calls where necessary. For `cmd/so-install/main.go`, we can extract logic to a testable function if necessary, or just test the wiring. For the TUI, we use `teatest` or manual model updates to reach 100% coverage.

### Risks
- Testing `os/user.Current()` and TUI logic may require careful mocking or extraction to prevent flaky tests.

### Ready for Proposal
Yes.