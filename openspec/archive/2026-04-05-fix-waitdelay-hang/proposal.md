# Proposal: Fix Installation Hang with Shorter WaitDelay

## Intent
Reduce the hang duration during software installation (specifically Ollama and OpenCode) caused by background daemons holding stdout/stderr pipes after the installation script exits.

## Scope
- `internal/infrastructure/shell/executor.go`: Separate `WaitDelay` from the main command timeout.
- `internal/infrastructure/shell/executor_test.go`: Add/update tests for `WaitDelay`.

## Approach
Currently, `ShellExecutor` sets `cmd.WaitDelay = e.timeout`. With a default timeout of 10 minutes, any daemonized process (like `ollama serve`) that inherits the installer script's pipes will cause `so-install` to hang for 10 minutes even after the script itself has exited normally.

The fix involves:
1. Adding a `waitDelay` field to `ShellExecutor` with a default value of 5 seconds.
2. Using this `waitDelay` for `cmd.WaitDelay` instead of the 10-minute timeout.
3. Updating the `NewShellExecutorWithTimeout` constructor to accept both durations for testing.

## Rollback Plan
Revert changes to `executor.go` and `executor_test.go`. The previous behavior (waiting 10 minutes) was safe but frustratingly slow.
