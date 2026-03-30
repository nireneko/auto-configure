# Tasks: System Update and Base Dependencies

## Phase 1: Domain / Foundation

- [x] 1.1 Add `SystemUpdate` and `BaseDeps` constants to `domain.SoftwareID` in `internal/core/domain/software.go`.
- [x] 1.2 Update `domain.GetSteps()` in `internal/core/domain/software.go` to include the "system-prep" step at index 0.
- [x] 1.3 Update `domain.AllSoftware()` in `internal/core/domain/software.go` (ensure new IDs are NOT included to keep them hidden from selection).
- [x] 1.4 Update `SoftwareID.DisplayName()` in `internal/core/domain/software.go` for the new IDs.

## Phase 2: Infrastructure Implementation (TDD)

- [x] 2.1 Create `internal/infrastructure/apt/update_test.go` with failing tests for `AptUpdateInstaller`.
- [x] 2.2 Create `internal/infrastructure/apt/update.go` and implement `AptUpdateInstaller` (must call `apt-get update` and `apt-get upgrade -y`).
- [x] 2.3 Create `internal/infrastructure/apt/deps_test.go` with failing tests for `BaseDepsInstaller`.
- [x] 2.4 Create `internal/infrastructure/apt/deps.go` and implement `BaseDepsInstaller` (must install `git`, `wget`, `curl`, etc.).

## Phase 3: Integration & Wiring

- [x] 3.1 Register `SystemUpdate` and `BaseDeps` installers in `cmd/so-install/main.go`.
- [x] 3.2 Update `internal/presentation/tui/model.go` to prepend `SystemUpdate` and `BaseDeps` to `m.selected` when "enter" is pressed in `stateSoftwareSelect`.

## Phase 4: Verification

- [x] 4.1 Run `make test` to verify all unit tests pass.
- [x] 4.2 Run `make lint` to ensure code quality.
- [x] 4.3 Verify the TUI shows "system-prep" phase during installation.
