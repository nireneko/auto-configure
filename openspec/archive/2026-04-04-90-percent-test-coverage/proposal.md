# Proposal: 90% Test Coverage and Edge Cases

## Intent
Increase project-wide test coverage to at least 90%, ensuring all packages meet this threshold and edge cases (like `SUDO_USER` and specific OS environments) are properly handled.

## Scope
- Refactor `cmd/so-install/main.go` to extract logic into a testable `Run` function.
- Increase coverage in `internal/infrastructure/nvidia`.
- Increase coverage in `internal/presentation/tui`.
- Increase coverage in `internal/infrastructure/osrelease`.
- Increase coverage in `internal/infrastructure/desktop`.
- Ensure `pkg/mocks` are exercised in tests.
- Add tests for edge cases identified during exploration.

## Approach
1. **Refactor `main.go`**: Move logic from `main()` to `Run(args []string, executor domain.Executor) int` to allow unit testing with mocks.
2. **Nvidia Coverage**: Add test cases for `installProprietaryDebian`, `installProprietaryNvidia`, and `enableNonFreeSources` by mocking shell output.
3. **TUI Coverage**: Expand `model_test.go` to cover `nextAfterNvidiaConfig`, `handleKey` (more keys), and all `view` methods.
4. **OSRelease Coverage**: Mock `/etc/os-release` and process lists to test `detectDesktopEnvironment` and `isProcessRunning` thoroughly.
5. **Desktop Coverage**: Add tests for KDE/Gnome specific logic in `screen_lock.go`.
6. **Mocks Coverage**: Ensure all mock methods in `pkg/mocks` are called at least once in tests that use them.

## Rollback Plan
- Changes are primarily in `*_test.go` files, which don't affect production.
- For `main.go`, the refactor will be verified by the new tests and manual execution. If it fails, revert to the previous `main.go` version.
