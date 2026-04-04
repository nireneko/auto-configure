# Tasks: Add Cursor and Google Antigravity IDEs

## Phase 1: Domain & Infrastructure Foundation

- [ ] 1.1 Add `Cursor` and `Antigravity` constants to `domain.SoftwareID` in `internal/core/domain/software.go`.
- [ ] 1.2 Update `domain.GetSteps()` to include `Cursor` and `Antigravity` in the `ides` step in `internal/core/domain/software.go`.
- [ ] 1.3 Update `domain.AllSoftware()` and `domain.DisplayName()` in `internal/core/domain/software.go`.

## Phase 2: Cursor Implementation

- [ ] 2.1 Create `internal/infrastructure/cursor/cursor.go` with `CursorInstaller` (download .deb, `apt install`).
- [ ] 2.2 Create `internal/infrastructure/cursor/cursor_test.go` and verify `IsInstalled` and `Install` logic with mocked executor.

## Phase 3: Antigravity Implementation

- [ ] 3.1 Create `internal/infrastructure/antigravity/antigravity.go` with `AntigravityInstaller` (add repo, `apt update`, `apt install`).
- [ ] 3.2 Create `internal/infrastructure/antigravity/antigravity_test.go` and verify `IsInstalled` and `Install` logic with mocked executor.

## Phase 4: Integration & Wiring

- [ ] 4.1 Update `cmd/so-install/main.go` to import the new packages and register the installers in `installerMap`.

## Phase 5: Verification

- [ ] 5.1 Run all tests in the project: `go test ./...`.
- [ ] 5.2 Build the project to ensure no compilation errors: `go build -o bin/so-install ./cmd/so-install`.
