# Design: Install Gentle-AI after AI tools

## Technical Approach
The implementation follows the established Clean Architecture pattern of the project. We will define the new software in the domain layer, implement a dedicated installer in the infrastructure layer, and register it in the main entry point. A new installation phase will be added to explicitly separate Gentle-AI from other AI tools as requested.

## Architecture Decisions

### Decision: Dedicated Installer Package
**Choice**: Create `internal/infrastructure/gentleai`.
**Alternatives considered**: Use a generic `ShellInstaller`.
**Rationale**: Existing tools like `ollama` and `opencode` have their own packages even if they just run shell commands. This maintains consistency and allows for tool-specific logic (like version parsing or post-install steps) in the future.

### Decision: Separate Installation Phase
**Choice**: New "gentle-ai" step in `GetSteps()`.
**Alternatives considered**: Add to "ai-cli" step.
**Rationale**: Fulfills the user requirement "once finished [with AI tools], the next thing that can be installed is Gentle-AI". It provides a clear visual transition in the TUI.

## Data Flow
The TUI model retrieves steps from `domain.GetSteps()`. When it reaches the "gentle-ai" step, it looks up the `GentleAIInstaller` in its map and calls `Install()`.

    TUI Model ──→ UseCase ──→ GentleAIInstaller ──→ ShellExecutor ──→ bash/sudo
                                     │
                                     └─→ GetActualUser() (from domain/user.go)

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/core/domain/software.go` | Modify | Add `GentleAI` ID, update `GetSteps` and `AllSoftware`. |
| `internal/infrastructure/gentleai/gentleai.go` | Create | Installer implementation using curl/bash script. |
| `internal/infrastructure/gentleai/gentleai_test.go` | Create | Unit tests for the installer using MockExecutor. |
| `cmd/so-install/main.go` | Modify | Register `GentleAI` installer in the global map. |

## Interfaces / Contracts

```go
// internal/infrastructure/gentleai/gentleai.go

type GentleAIInstaller struct {
    executor domain.Executor
    userName string
}

func NewGentleAIInstaller(executor domain.Executor) *GentleAIInstaller {
    return &GentleAIInstaller{
        executor: executor,
        userName: domain.GetActualUser(),
    }
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `IsInstalled` | Mock executor to return success/failure for `gentle-ai --version`. |
| Unit | `Install` | Mock executor to verify the `curl | bash` command is called correctly. |
| Unit | User Context | Verify `sudo -u` is used when `SUDO_USER` is set. |

## Migration / Rollout
No migration required. This is a new optional feature.
