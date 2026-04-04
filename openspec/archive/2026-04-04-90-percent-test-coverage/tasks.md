# Tasks: 90% Test Coverage

## Phase 1: Entry Point & Mocks (Target: 90%)
- [x] 1.1 Refactor `cmd/so-install/main.go`: extract logic to `Run`.
- [x] 1.2 Create `cmd/so-install/main_test.go`.
- [x] 1.3 Add tests for `pkg/mocks/software_installer.go` (ensure `Install` and `IsInstalled` are called).

## Phase 2: Infrastructure Layer (Target: 90%)
- [x] 2.1 Update `internal/infrastructure/nvidia/nvidia_test.go` to cover proprietary/free branches.
- [x] 2.2 Update `internal/infrastructure/osrelease/detector_test.go` for `detectDesktopEnvironment`.
- [x] 2.3 Update `internal/infrastructure/desktop/screen_lock_test.go` for KDE/Gnome logic.
- [x] 2.4 Update `internal/infrastructure/homebrew/homebrew_test.go` for shell configuration logic.
- [ ] 2.5 Update `internal/infrastructure/gitlab/configurator_test.go` for all config branches.

## Phase 3: Presentation Layer (Target: 90%)
- [x] 3.1 Update `internal/presentation/tui/model_test.go` to cover `nextAfterNvidiaConfig`.
- [x] 3.2 Add tests for all `view*` methods in `model.go`.
- [x] 3.3 Expand `handleKey` tests to cover edge cases and all navigation options.

## Phase 4: Verification
- [ ] 4.1 Run `go test -cover ./...`.
- [ ] 4.2 Verify all packages are >= 90%.
