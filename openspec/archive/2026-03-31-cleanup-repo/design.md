# Design: Repository Cleanup and .gitignore

## Technical Approach
Implement a standardized `.gitignore` for Go projects and remove specifically identified junk files that are currently cluttering the repository root or being tracked by Git.

## Architecture Decisions

### Decision: .gitignore Content
**Choice**: Use a comprehensive Go-standard `.gitignore` template.
**Alternatives considered**: Minimalist manually maintained list.
**Rationale**: Prevents accidental tracking of common temporary files, IDE configs, and OS junk without needing constant manual updates.

### Decision: Handling of Existing Tracked Junk
**Choice**: Use `git rm --cached` followed by physical deletion.
**Alternatives considered**: Only `git rm` (deletes and stops tracking in one go), or just `rm` (leaves it as "deleted" in git status).
**Rationale**: `git rm` is the standard way to stop tracking and delete. `git rm --cached` is useful if we want to keep it locally but stop tracking, but here we want them GONE.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `.gitignore` | Modify | Expand with Go, IDE, and OS ignore patterns. |
| `so-install` | Delete | Root binary artifact (already exists in `bin/` if built correctly). |
| `error-ddev.png` | Delete | Non-essential screenshot. |
| `coverage.out` | Delete | Temporary test output. |

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Manual | Git status | Verify no untracked or tracked junk remains. |
| Manual | File existence | Verify `error-ddev.png` and root `so-install` are physically gone. |
| Manual | Ignore verification | Create a dummy `.exe` or `vendor/` and ensure it's ignored. |

## Migration / Rollout
No migration required. This is a one-time cleanup.

## Open Questions
None.
