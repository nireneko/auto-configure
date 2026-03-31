# Proposal: Repository Cleanup and .gitignore

## Intent
Clean the repository of binary artifacts, screenshots, and unnecessary Go files that are currently tracked or present in the root directory. Establish a comprehensive `.gitignore` to prevent future pollution.

## Scope

### In Scope
- Create a robust `.gitignore` for a Go project.
- Remove tracked binary `so-install` from root.
- Remove tracked screenshot `error-ddev.png`.
- Remove untracked junk like `coverage.out` and root `bin/` contents.

### Out of Scope
- Modifying `go.mod` or `go.sum`.
- Modifying source code in `cmd/`, `internal/`, or `pkg/`.

## Approach
1. Update `.gitignore` with standard Go ignore patterns.
2. Remove `so-install` and `error-ddev.png` from git tracking (keep local or delete as requested).
3. Physically delete the junk files.
4. Verify repo status with `git status`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `.gitignore` | Modified | Added comprehensive Go ignore patterns. |
| `so-install` | Removed | Tracked binary in root. |
| `error-ddev.png` | Removed | Tracked screenshot. |
| `coverage.out` | Removed | Untracked junk. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Deleting source code | Low | Double-check file list before `rm`. |
| Ignoring necessary files | Low | Use standard Go `.gitignore` patterns. |

## Rollback Plan
- Use `git checkout .` to restore tracked files.
- Manual recreation of `.gitignore` if needed.

## Dependencies
- None.

## Success Criteria
- [ ] `error-ddev.png` is gone from filesystem and git.
- [ ] `so-install` is gone from root filesystem and git.
- [ ] `coverage.out` is gone.
- [ ] `.gitignore` contains common Go patterns.
- [ ] `git status` shows a clean working tree (no untracked junk).
