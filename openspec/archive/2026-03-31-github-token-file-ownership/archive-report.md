# Archive Report: github-token-file-ownership

**Archived**: 2026-03-31
**Change**: Fix file ownership for GitlabTokenConfigurator under sudo
**Status**: COMPLETE
**All Artifacts**: Engram (persistent memory)

---

## Executive Summary

The `github-token-file-ownership` change successfully fixed a critical issue where `GitlabTokenConfigurator.Install()` running via `sudo` created `~/.composer/auth.json` and `~/.npmrc` owned by root, making them unreadable/unwritable by the real user. The solution added `GetActualUID()` and `GetActualGID()` domain helpers to read `SUDO_UID`/`SUDO_GID` env vars, and injected `chownFn` into the configurator to chown files back to the real user after every write. All 10 tasks completed in Strict TDD mode (RED → GREEN → VERIFY). Zero test failures; full suite passing.

---

## Artifacts

| Artifact | Observation ID | Status |
|----------|---|---|
| Proposal | #518 | ✅ Complete |
| Specification | #519 | ✅ Complete |
| Technical Design | #520 | ✅ Complete |
| Task Breakdown | #521 | ✅ Complete (10/10 tasks) |
| Apply Progress | #522 | ✅ Complete |
| Verification Report | #523 | ✅ PASS (all requirements met) |

**Persistence**: All artifacts stored in Engram with topic keys `sdd/github-token-file-ownership/{artifact-type}` for cross-session traceability.

---

## Specs Synced to Main Specs

| Domain | Changes | Details |
|--------|---------|---------|
| `openspec/specs/domain/spec.md` | ADDED 2 requirements | GetActualUID, GetActualGID (3 scenarios each) |
| `openspec/specs/infrastructure/spec.md` | ADDED 3 requirements | Composer chown, NPM chown, chownFn injection (6 scenarios total) |

**Merge approach**: Added new sections to existing main specs; preserved all prior requirements.

---

## Implementation Summary

### Code Changes

**Files modified** (4):
- `internal/core/domain/user.go` — Added `GetActualUID()` and `GetActualGID()`
- `internal/core/domain/user_test.go` — Added `TestGetActualUID` and `TestGetActualGID` (6 sub-tests total)
- `internal/infrastructure/gitlab/configurator.go` — Added `chownFn`, `uidFn`, `gidFn` fields; setters; 3 chown call sites
- `internal/infrastructure/gitlab/configurator_test.go` — Added chown verification tests (2 new sub-tests)

### Key Implementation Details

**Domain helpers** (`user.go`):
```go
GetActualUID() int {
    if s := os.Getenv("SUDO_UID"); s != "" {
        if uid, err := strconv.Atoi(s); err == nil {
            return uid
        }
    }
    return os.Getuid()
}

GetActualGID() int {
    if s := os.Getenv("SUDO_GID"); s != "" {
        if gid, err := strconv.Atoi(s); err == nil {
            return gid
        }
    }
    return os.Getgid()
}
```

**Dependency injection** (`configurator.go`):
- Added struct fields: `chownFn func(string, int, int) error`, `uidFn func() int`, `gidFn func() int`
- Constructor defaults: `chownFn: os.Chown`, `uidFn: domain.GetActualUID`, `gidFn: domain.GetActualGID`
- Setters: `SetChownFn`, `SetUIDFn`, `SetGIDFn` for test injection

**Chown call sites** (3 total):
1. After `os.MkdirAll(composerDir)` in `configureComposer()`
2. After `os.WriteFile(authFile)` in `configureComposer()`
3. After `os.WriteFile(npmrcFile)` in `configureNpm()`

All calls wrapped with error messages: `"failed to chown <resource>"`

---

## Verification Results

**Test Execution**:
- `go test ./internal/core/domain/...` — ✅ PASS (6 new tests + all existing)
- `go test ./internal/infrastructure/gitlab/...` — ✅ PASS (2 new tests + all existing)
- `go test ./...` (full suite) — ✅ PASS (zero failures, zero regressions)

**Spec Compliance**:
| Requirement | Scenario Count | Passing | Status |
|---|---|---|---|
| GetActualUID | 3 | 3 | ✅ COMPLIANT |
| GetActualGID | 3 | 3 | ✅ COMPLIANT |
| Composer chown under sudo | 2 | 2 | ✅ COMPLIANT |
| NPM chown under sudo | 2 | 2 | ✅ COMPLIANT |
| chownFn injection | 2 | 2 | ✅ COMPLIANT |
| **TOTAL** | **12** | **12** | **✅ 100% COMPLIANT** |

**Correctness Checks**:
- GetActualUID/GID parse and fallback logic verified
- All 3 chown call sites present and error-wrapped
- Struct fields initialized to correct defaults
- All 3 setter methods present
- Test suite uses spy injection pattern correctly

**Verdict**: PASS — All requirements met, all tests passing, zero regressions.

---

## Design Rationale

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Injection pattern | Struct fields + setters | Matches existing patterns in codebase (`SetHomeDir`, `SetUserFn` in `screen_lock.go`) |
| UID/GID source | `SUDO_UID`/`SUDO_GID` env vars | Simpler than `user.Lookup(SUDO_USER)`; always set by sudo; avoids string-to-int conversion |
| Fallback | `os.Getuid()`/`os.Getgid()` | Safe default; chowning to own UID is a no-op on Linux |
| Function fields | `uidFn`, `gidFn` injectable | Enables clean test injection without env var manipulation |

---

## Test Coverage

**New tests**:
- `TestGetActualUID`: SUDO_UID set → returns 1000; unset → returns os.Getuid(); invalid → fallback
- `TestGetActualGID`: SUDO_GID set → returns 1000; unset → returns os.Getgid(); invalid → fallback
- Configurator chown: Spy verifies 3 calls with correct paths (composerDir, auth.json, .npmrc) and uid/gid=1000
- Configurator error propagation: chownFn error bubbles up through Install()

**Existing test behavior**:
- All pre-existing tests continued passing because `os.Chown` is a no-op when chowning to the file's current owner
- No test modifications required; additive only

---

## Risks & Mitigations

| Risk | Likelihood | Mitigation | Status |
|------|---|---|---|
| `SUDO_UID` not set on some distros | Low | Fallback to `os.Getuid()` (no-op chown) | ✅ Tested |
| Chown fails on non-root execution | Low | `os.Chown` to own UID is no-op; error gracefully wrapped | ✅ Tested |
| Race condition between write and chown | Very Low | File created immediately; chown follows synchronously | ✅ Inherent safety |

**Status**: No residual risks. All scenarios covered by tests.

---

## Rollback Plan

Revert commits affecting:
- `internal/core/domain/user.go` and `user_test.go`
- `internal/infrastructure/gitlab/configurator.go` and `configurator_test.go`

No database, config, or external state affected. Behavior is purely additive.

---

## SDD Cycle Status

| Phase | Status | Tasks | Result |
|---|---|---|---|
| Explore | ✅ Done | — | Identified solution |
| Propose | ✅ Done | — | Defined scope and approach |
| Spec | ✅ Done | — | 12 requirements + scenarios |
| Design | ✅ Done | — | Architecture & implementation details |
| Tasks | ✅ Done | 10 | All completed |
| Apply | ✅ Done | 10/10 | Strict TDD: RED → GREEN → VERIFY |
| Verify | ✅ Done | — | All 12 scenarios passing, zero regressions |
| **Archive** | ✅ Done | — | **CYCLE COMPLETE** |

---

## Next Steps

None. Change is fully archived and ready for integration. The fix is production-ready.

---

## Traceability

All SDD artifacts live in Engram for cross-session recovery:
- `sdd/github-token-file-ownership/proposal` (#518)
- `sdd/github-token-file-ownership/spec` (#519)
- `sdd/github-token-file-ownership/design` (#520)
- `sdd/github-token-file-ownership/tasks` (#521)
- `sdd/github-token-file-ownership/apply-progress` (#522)
- `sdd/github-token-file-ownership/verify-report` (#523)
- `sdd/github-token-file-ownership/archive-report` (this document, saved to Engram)

Spec syncs reflected in:
- `openspec/specs/domain/spec.md` (GetActualUID, GetActualGID)
- `openspec/specs/infrastructure/spec.md` (Composer/NPM chown, chownFn injection)
