# Tasks: Fix Ollama & OpenCode Installation Hang

## Phase 1: Implementation

- [x] 1.1 In `internal/infrastructure/shell/executor.go`: add `timeout time.Duration` field to `ShellExecutor` struct and set it to `10 * time.Minute` in `NewShellExecutor()`
- [x] 1.2 Add exported `NewShellExecutorWithTimeout(d time.Duration) *ShellExecutor` constructor for test injection
- [x] 1.3 Rewrite `Execute()` to use `context.WithTimeout(context.Background(), e.timeout)` + `exec.CommandContext(ctx, name, args...)` + `defer cancel()`
- [x] 1.4 Add `cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}` before `cmd.Run()` in `Execute()`

## Phase 2: Testing (TDD — RED → GREEN)

- [x] 2.1 **RED**: In `executor_test.go`, add `TestShellExecutor_Timeout` — use `NewShellExecutorWithTimeout(500ms)` + `Execute("sh", "-c", "sleep 999")` — assert non-nil error returned within ~1s
- [x] 2.2 **RED**: Add `TestShellExecutor_DaemonDoesNotHang` — use `NewShellExecutorWithTimeout(2s)` + `Execute("sh", "-c", "sleep 999 &")` — assert Execute returns before timeout fires (i.e. the background sleep doesn't block)
- [x] 2.3 **GREEN**: Verify both new tests pass with the implementation from Phase 1
- [x] 2.4 Confirm existing tests (`TestShellExecutor_SuccessfulCommand`, `TestShellExecutor_FailingCommand`, `TestShellExecutor_CapturesStderr`) still pass unchanged

## Phase 3: Verification

- [x] 3.1 Run `go test ./internal/infrastructure/shell/...` — all tests must pass
- [x] 3.2 Run `go test ./...` — no regressions across the full suite
- [x] 3.3 Run `go vet ./...` — no issues
