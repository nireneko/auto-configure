# Tasks: Install Gentle-AI after AI tools

## Phase 1: Domain / Foundation

- [x] 1.1 Add `GentleAI SoftwareID = "gentle-ai"` to `internal/core/domain/software.go`.
- [x] 1.2 Update `DisplayName()` in `internal/core/domain/software.go` to return "Gentle-AI" for the new ID.
- [x] 1.3 Update `AllSoftware()` in `internal/core/domain/software.go` to include `GentleAI`.
- [x] 1.4 Update `GetSteps()` in `internal/core/domain/software.go` to insert the "gentle-ai" step after "ai-cli".

## Phase 2: Infrastructure Implementation

- [x] 2.1 Create directory `internal/infrastructure/gentleai/`.
- [x] 2.2 Implement `GentleAIInstaller` in `internal/infrastructure/gentleai/gentleai.go`.
- [x] 2.3 Implement `IsInstalled()` using `gentle-ai --version` check.
- [x] 2.4 Implement `Install()` using the official curl-based installation script.
- [x] 2.5 Ensure `Install()` respects the actual user context (sudo -u) using `domain.GetActualUser()`.

## Phase 3: Testing

- [x] 3.1 Create `internal/infrastructure/gentleai/gentleai_test.go`.
- [x] 3.2 Add tests for `IsInstalled` covering both installed and not installed states.
- [x] 3.3 Add tests for `Install` verifying the correct command string and sudo wrapping.
- [x] 3.4 Verify 100% test coverage for the new package.

## Phase 4: Integration & Wiring

- [x] 4.1 Register `GentleAI` installer in `cmd/so-install/main.go`.
- [x] 4.2 Run `make test` to ensure no regressions in domain or main packages.
- [x] 4.3 Manually verify that "Gentle-AI" appears in the TUI selection list.
