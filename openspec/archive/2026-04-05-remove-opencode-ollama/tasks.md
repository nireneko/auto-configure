# Tasks: Remove OpenCode and Ollama

## Phase 1: Domain & Entry Point
- [x] 1.1 Remove `Ollama` and `OpenCode` constants from `internal/core/domain/software.go`
- [x] 1.2 Update `GetSteps()`, `AllSoftware()` and `DisplayName()` in `software.go`
- [x] 1.3 Fix `cmd/so-install/main.go` (imports and installerMap)

## Phase 2: Infrastructure & Cleanup
- [x] 2.1 Delete `internal/infrastructure/opencode/`
- [x] 2.2 Delete `internal/infrastructure/ollama/`
- [x] 2.3 Fix `internal/core/domain/software_test.go`

## Phase 3: Documentation & Specs
- [x] 3.1 Update `README.md` and `README_es.md`
- [x] 3.2 Update `openspec/specs/repository/spec.md`
- [x] 3.3 Update `openspec/specs/infrastructure/spec.md`

## Phase 4: Final Verification
- [x] 4.1 Run `make test` and ensure success
- [x] 4.2 Create verification report
