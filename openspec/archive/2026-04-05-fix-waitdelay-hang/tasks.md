# Tasks: Fix Installation Hang with Shorter WaitDelay

## Phase 1: Infrastructure (Shell Executor)
- [x] 1.1 Update `ShellExecutor` struct with `waitDelay` field
- [x] 1.2 Update `NewShellExecutor` and `NewShellExecutorWithTimeout` constructors
- [x] 1.3 Update `Execute` method to use `cmd.WaitDelay = e.waitDelay`

## Phase 2: Implementation (Verification)
- [x] 2.1 Update `executor_test.go` to test isolated `WaitDelay` behavior
- [x] 2.2 Ensure existing tests pass (`make test`)

## Phase 3: Archive
- [x] 3.1 Create verification report
- [x] 3.2 Sync specs to main repository
- [x] 3.3 Archive the change
