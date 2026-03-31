# Tasks: Gitlab Token Configuration

## Phase 1: Domain & Infrastructure Foundation

- [x] 1.1 Add `GitlabTokenConfig` to `domain.SoftwareID` in `internal/core/domain/software.go`.
- [x] 1.2 Update `domain.GetSteps()` to include `GitlabTokenConfig` in the `apps` step.
- [x] 1.3 Create `internal/infrastructure/gitlab/configurator.go` with `GitlabTokenConfigurator` struct.
- [x] 1.4 Implement `SoftwareInstaller` interface for `GitlabTokenConfigurator`.
- [x] 1.5 Add `SetToken(string)` method to `GitlabTokenConfigurator`.
- [x] 1.6 Implement `Install()` in `GitlabTokenConfigurator` to update `~/.composer/auth.json` and `~/.npmrc`.

## Phase 2: TUI Implementation

- [x] 2.1 Add `stateTokenInput` and `gitlabToken` field to `Model` in `internal/presentation/tui/model.go`.
- [x] 2.2 Update `handleKey` in `model.go` to manage the transition from `stateSoftwareSelect` to `stateTokenInput`.
- [x] 2.3 Implement token input handling in `handleKey` for `stateTokenInput` (capture chars, backspace, enter).
- [x] 2.4 Add `viewTokenInput` to `internal/presentation/tui/model.go` (or a new file) with masked input display.
- [x] 2.5 Update `View()` in `model.go` to include the new `viewTokenInput`.

## Phase 3: Wiring & Integration

- [x] 3.1 Update `cmd/so-install/main.go` to instantiate `GitlabTokenConfigurator`.
- [x] 3.2 Update `Update()` in `model.go` to pass the captured token to the configurator before `startInstallMsg`.
- [x] 3.3 Ensure the token is passed correctly through the `InstallSoftwareUseCase`.

## Phase 4: Verification & Testing

- [x] 4.1 Create `internal/infrastructure/gitlab/configurator_test.go` and verify file updates.
- [x] 4.2 Add unit tests for TUI state transition in `internal/presentation/tui/model_test.go`.
- [x] 4.3 Verify that `~/.composer/auth.json` and `~/.npmrc` are created/updated correctly with a mock home directory.
