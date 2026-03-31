# Tasks: 100% Test Coverage

## Phase 1: Test Infrastructure

- [ ] 1.1 Create `pkg/mocks/software_installer_test.go` to test `ID()` and mock instantiations.

## Phase 2: Domain Layer

- [ ] 2.1 Create `internal/core/domain/errors_test.go` to cover all `Error()` string outputs and `WrapInstallError()`.
- [ ] 2.2 Create `internal/core/domain/user_test.go` to cover `GetActualUser()` and `GetActualHome()` using `t.Setenv("SUDO_USER", "mockuser")`.

## Phase 3: Infrastructure Layer

- [ ] 3.1 Modify `internal/infrastructure/osrelease/detector_test.go` to test `detectDesktopEnvironment` logic and `isProcessRunning`.
- [ ] 3.2 Modify `internal/infrastructure/gitlab/configurator_test.go` to cover PHP Composer and NPM package configurations.
- [ ] 3.3 Modify all infrastructure installers (e.g. `browsers`, `ddev`, `flatpak`, `openvpn`, `desktop`, `docker`, `homebrew`, `npm`, `nvm`) to test their `ID()` method and any uncovered branches.

## Phase 4: Presentation Layer

- [ ] 4.1 Modify `internal/presentation/tui/model_test.go` to manually send `tea.KeyMsg` to `Update()` and test transitions to `stateSoftwareSelect`, `stateTokenInput`, `stateProgress`, and `stateSummary`.
- [ ] 4.2 Test the `View()` function for all states in `model_test.go`.

## Phase 5: Main Entry Point

- [ ] 5.1 Refactor `cmd/so-install/main.go` to extract the main logic into a `run(args []string) error` function.
- [ ] 5.2 Create `cmd/so-install/main_test.go` to test the `run()` function using mock arguments or environment variables.

## Phase 6: Verification

- [ ] 6.1 Run `make test` or `go test -cover ./...` and ensure 100.0% coverage across all packages.