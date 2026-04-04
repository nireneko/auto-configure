# Design: Translate README to English

## Technical Approach
We will produce a high-quality translation of the Spanish README (`README_es.md`) into a new `README.md` file in English. The original Spanish version will be kept as `README_es.md`.

## Decisions

### Decision: Translation Language
**Choice**: English.
**Rationale**: Direct user request.

### Decision: Preservation of Code Blocks
**Choice**: Preserve all shell command blocks and Makefile targets.
**Rationale**: Commands like `make build` and `make run` are identical in any language.

## Content Outline (English)
- **1X-SO Install Orchestrator**: Title and high-level description.
- **Main Features**: Deterministic post-installation tool.
- **What it Installs and Configures**: Categorized list of modules.
- **System Requirements**: Go 1.24.2+, Make, Debian-based OS, Sudo.
- **Usage**:
    - Building: `make build`
    - Running: `make run`
    - Testing: `make test`

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `README.md` | Create | English version of the project documentation. |
| `README_es.md` | Rename | Kept from Spanish version for reference. |
