# Proposal: Fix Ollama & OpenCode Installation Hang

## Intent

The installation of Ollama and OpenCode freezes indefinitely. Root cause: `ShellExecutor.Execute()` uses `cmd.Run()` with `bytes.Buffer` for stdout/stderr — no timeout, no context. The install scripts spawn daemon processes (via fork-without-exec) that inherit the pipe write-end, preventing EOF and causing `cmd.Run()` to block forever.

## Scope

### In Scope
- Add `context.WithTimeout` to `ShellExecutor.Execute()`
- Set `syscall.SysProcAttr{Setpgid: true}` to isolate the process group
- Update `ShellExecutor` tests to cover timeout and Setpgid behavior
- Update `domain.Executor` interface if needed

### Out of Scope
- Fixing OpenCode's potential PATH issue with sudo and NVM (separate concern)
- Per-installer timeout configuration
- Streaming output to TUI during installation

## Approach

Add a 10-minute `context.WithTimeout` to `ShellExecutor.Execute()`. Set `Setpgid: true` so the child process runs in its own process group — when the context cancels and sends SIGKILL, the entire group is killed, closing all inherited pipe FDs and unblocking `cmd.Wait()`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/infrastructure/shell/executor.go` | Modified | Add context, timeout, SysProcAttr |
| `internal/infrastructure/shell/executor_test.go` | Modified | Add timeout and Setpgid tests |
| `internal/core/domain/executor.go` | None | Interface unchanged (`Execute(name, args...)`) |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Timeout too short for slow systems | Low | 10 min is generous for any script-based installer |
| `Setpgid` doesn't kill grandchildren that call `setsid()` | Low | Ollama/OpenCode scripts don't do this |
| Breaking existing tests that mock `Executor` | Low | Interface signature unchanged |

## Rollback Plan

Revert `internal/infrastructure/shell/executor.go` to the version without context/SysProcAttr. The interface and all callers remain unchanged.

## Dependencies

- `syscall` package (stdlib, already available on Linux)
- `context` package (stdlib)

## Success Criteria

- [ ] Ollama installation completes without hanging
- [ ] OpenCode installation completes without hanging
- [ ] `ShellExecutor` tests pass including new timeout/Setpgid coverage
- [ ] `go test ./...` passes with no regressions
