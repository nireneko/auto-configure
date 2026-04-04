# Tasks: Install VS Code

## Phase 1: Domain & Tests (Foundation)

- [x] 1.1 Update `internal/core/domain/software.go`: Add `VsCode` SoftwareID constant.
- [x] 1.2 Update `internal/core/domain/software.go`: Add `ides` step in `GetSteps()` and update `AllSoftware()`.
- [x] 1.3 Update `internal/core/domain/software.go`: Update `DisplayName()` to return "Visual Studio Code" for `VsCode`.
- [x] 1.4 Update `internal/core/domain/software_test.go`: Add test cases for `VsCode` display name and `ides` step inclusion.

## Phase 2: Infrastructure (Implementation)

- [x] 2.1 Create `internal/infrastructure/vscode/vscode.go`: Implement `VsCodeInstaller` with `Install()`, `IsInstalled()`, and `ID()`.
- [x] 2.2 Create `internal/infrastructure/vscode/vscode_test.go`: Write unit tests for `VsCodeInstaller` using `MockExecutor`.
- [x] 2.3 Verify `VsCodeInstaller.Install` calls `wget` with the correct URL and `apt install` with the downloaded `.deb`.
- [x] 2.4 Verify `VsCodeInstaller.IsInstalled` calls `which code`.

## Phase 3: Wiring (Integration)

- [x] 3.1 Update `cmd/so-install/main.go`: Register `VsCodeInstaller` in the `installerMap`.
- [x] 3.2 Update `PRD.md`: Add a new section for IDEs (VS Code) in Functional Requirements.

## Phase 4: Verification & Cleanup

- [x] 4.1 Run all tests: `go test ./...`.
- [ ] 4.2 Verify TUI: Run `go run cmd/so-install/main.go` (dry run or mock if possible) to see VS Code in the list.
- [x] 4.3 Documentation: Ensure all new code follows project standards and comments are present where necessary.
