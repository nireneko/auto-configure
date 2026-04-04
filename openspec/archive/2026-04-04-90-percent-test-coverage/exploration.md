## Exploration: 90% Test Coverage and Edge Cases

### Current State
The project has an overall coverage of 87.9%. Several packages are below the 90% target:
- `cmd/so-install`: 50.0%
- `internal/infrastructure/nvidia`: 77.3%
- `pkg/mocks`: 80.0%
- `internal/infrastructure/desktop`: 82.8%
- `internal/presentation/tui`: 84.3%
- `internal/infrastructure/osrelease`: 86.9%
- `internal/infrastructure/homebrew`: 88.6%
- `internal/infrastructure/gitlab`: 88.7%
- `internal/infrastructure/nvm`: 89.1%

### Affected Areas
- `cmd/so-install/main.go` — `main` function is partially covered (52.0%). Needs refactoring to be testable.
- `internal/infrastructure/nvidia/nvidia.go` — Proprietary install methods and kernel headers are partially covered (61-66%).
- `internal/presentation/tui/model.go` — Several methods like `nextAfterNvidiaConfig` (40.0%), `handleKey` (81.2%), and various `view` methods are below 90%.
- `internal/infrastructure/osrelease/detector.go` — `detectDesktopEnvironment` (33.3%) and `NewDefaultDetector` (50.0%) need better coverage.
- `internal/infrastructure/desktop/screen_lock.go` — KDE/Gnome specific installation logic is partially covered.
- `pkg/mocks/software_installer.go` — Mock methods `Install` and `IsInstalled` have 0% coverage.

### Approaches
1. **Incremental Coverage Improvement** — Target the specific files and functions identified in the `go tool cover` output.
   - Pros: Surgical, minimal changes to working code (except for `main.go`).
   - Cons: Might miss broader architectural improvements for testability.
   - Effort: Medium

2. **Refactor for Testability + Comprehensive Mocking** — Refactor `main.go` and improve mock usage across the board.
   - Pros: Higher quality code, better long-term maintainability, easier to hit 100% later.
   - Cons: More intrusive changes to `main.go`.
   - Effort: Medium-High

### Recommendation
Approach 1 combined with the refactor of `main.go` (from Approach 2) is recommended. This allows hitting the 90% target efficiently while fixing the most "untestable" part of the codebase.

### Risks
- `main.go` refactoring might introduce regressions if not handled carefully (though tests will be added).
- TUI testing with `bubbletea` can be tricky for complex state transitions.
- Hardware-specific logic (Nvidia) requires careful mocking of shell commands.

### Ready for Proposal
Yes — I have a clear map of what's missing and how to fix it.
