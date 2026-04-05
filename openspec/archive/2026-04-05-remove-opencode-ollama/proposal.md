# Proposal: Remove OpenCode and Ollama

## Intent
Remove OpenCode and Ollama from the project as requested by the user.

## Scope
- `internal/core/domain/software.go`: Constants and step definitions.
- `cmd/so-install/main.go`: Installer map.
- `internal/infrastructure/opencode/`: Implementation and tests.
- `internal/infrastructure/ollama/`: Implementation and tests.
- `README.md`, `README_es.md`: Documentation.
- `openspec/specs/`: Specifications.

## Approach
Completely remove all references, code, and tests related to `OpenCode` and `Ollama`. Ensure the project builds and all remaining tests pass.
