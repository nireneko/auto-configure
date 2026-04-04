# Design: Fix Ollama & OpenCode Installation Hang

## Technical Approach

Modify `ShellExecutor.Execute()` to wrap every command in a `context.WithTimeout(10 minutes)` and set `cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}`. No changes to the `Executor` interface or any caller. Add two new test cases in `executor_test.go`.

## Architecture Decisions

### Decision: Timeout location — executor vs. installer

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Timeout in `ShellExecutor` (global) | All commands get a ceiling; simple, one place | ✅ Chosen |
| Timeout per-installer (e.g. `OllamaInstaller`) | Granular, but requires interface change + duplication | ❌ Rejected |
| Context passed by caller | Flexible, but breaks `domain.Executor` interface | ❌ Rejected |

**Rationale**: The `domain.Executor` interface must remain stable. A 10-minute global ceiling is safe for all known installers and requires zero changes to callers.

### Decision: Setpgid to break pipe inheritance

| Option | Tradeoff | Decision |
|--------|----------|----------|
| `Setpgid: true` | Child runs in own process group; context SIGKILL kills the group | ✅ Chosen |
| `Pdeathsig: SIGKILL` | Only kills direct child, not grandchildren | ❌ Insufficient |
| Redirect stdout to `os.DevNull` | Eliminates pipe, but loses stderr for error reporting | ❌ Rejected |

**Rationale**: `Setpgid: true` combined with `exec.CommandContext` is the minimal change that ensures the timeout fires and unblocks `cmd.Wait()` even when daemons inherit pipe FDs.

## Data Flow

```
ShellExecutor.Execute(name, args)
    │
    ├── context.WithTimeout(10 min)
    ├── exec.CommandContext(ctx, name, args...)
    ├── cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    ├── cmd.Stdout = &outBuf  (bytes.Buffer, unchanged)
    ├── cmd.Stderr = &errBuf  (bytes.Buffer, unchanged)
    └── cmd.Run()
            │
            ├── script exits normally → returns stdout, stderr, nil
            └── daemon inherits pipe → timeout fires → SIGKILL to group
                                     → pipe closes → cmd.Run() returns ctx.Err()
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/infrastructure/shell/executor.go` | Modify | Add context, timeout, SysProcAttr |
| `internal/infrastructure/shell/executor_test.go` | Modify | Add timeout and hang-prevention tests |

## Interfaces / Contracts

`domain.Executor` interface is **unchanged**:
```go
Execute(name string, args ...string) (stdout, stderr string, err error)
```

`ShellExecutor.Execute` internal signature after change:
```go
func (e *ShellExecutor) Execute(name string, args ...string) (string, string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()
    cmd := exec.CommandContext(ctx, name, args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    // ... rest unchanged
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Command exits normally → stdout/stderr captured | existing tests (unchanged) |
| Unit | Command fails → non-nil error | existing test (unchanged) |
| Unit | stderr captured on failure | existing test (unchanged) |
| Unit | Command exceeds timeout → returns error | new: `sleep 999`, short timeout via a testable executor variant |
| Unit | Long-running bg daemon doesn't block | new: script that forks background sleep, verify Execute returns |

> **Note**: The timeout test requires a way to inject a short timeout. Two approaches:
> - Export a `NewShellExecutorWithTimeout(d time.Duration)` constructor for testing only.
> - Use `sh -c "sleep 999"` with a 1-second timeout via the testable constructor.

## Migration / Rollout

No migration required. No interface changes. All existing callers and mocks are unaffected.

## Open Questions

- None.
