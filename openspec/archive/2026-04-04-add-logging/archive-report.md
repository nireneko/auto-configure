# Change Report: Add Logging for Debugging and Visibility

## Goal
Implement a structured logging system to record all critical operations, including shell command execution, TUI state transitions, and errors, to help diagnose why the application might appear to freeze.

## Accomplished
- ✅ Defined `domain.Logger` interface to decouple logging from core logic.
- ✅ Implemented `logging.FileLogger` in the infrastructure layer using `log/slog`.
- ✅ Injected the logger into `ShellExecutor` to log command execution details (stdout, stderr, errors, timeouts).
- ✅ Injected the logger into the TUI `Model` to log state transitions and installation progress.
- ✅ Initialized the logger in `main.go` to write to `so-install.log` in the current directory.
- ✅ Updated all existing tests to support the modified constructors.
- ✅ Added documentation about logs to `README.md` and `README_es.md`.

## Verified
- Ran all tests using `make test`: All passed.
- Verified compilation with `make build`: Successful.

## Next Steps
- (Optional) Implement different log levels (Debug, Info, Error) via CLI flags.
- (Optional) Add more granular logging in complex use cases or individual installers.

## Relevant Files
- `internal/core/domain/logger.go`
- `internal/infrastructure/logging/logger.go`
- `internal/infrastructure/shell/executor.go`
- `internal/presentation/tui/model.go`
- `cmd/so-install/main.go`
- `README.md`
- `README_es.md`
