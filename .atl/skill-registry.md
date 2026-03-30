# Skill Registry - 1x-so-install

## Compact Rules

### Go Testing
- Use `go test` for running tests.
- Use `go test -cover` for coverage.
- Mocking is preferred for infrastructure components (see `pkg/mocks`).
- Follow the patterns in `internal/core/usecases/*_test.go` and `internal/infrastructure/*/*_test.go`.

### Clean Architecture
- Business logic in `internal/core/usecases`.
- Domain entities and interfaces in `internal/core/domain`.
- Implementation details (shell commands, file access) in `internal/infrastructure`.
- TUI logic in `internal/presentation/tui`.

## User Skills

| Skill | Trigger |
|-------|---------|
| go-testing | Go tests, Bubbletea TUI testing |
| skill-creator | Creating new AI skills |
| branch-pr | Creating pull requests |
| issue-creation | Creating issues or reporting bugs |
| judgment-day | Reviewing code or dual review |
| sdd-init | Initializing SDD context |
| sdd-explore | Investigating ideas |
| sdd-propose | Proposing changes |
| sdd-spec | Writing specifications |
| sdd-design | Creating technical designs |
| sdd-tasks | Creating implementation tasks |
| sdd-apply | Implementing changes |
| sdd-verify | Verifying implementation |
| sdd-archive | Archiving changes |
