# Design: 90% Test Coverage and Edge Cases

## Architecture Decisions

### Decision: Main Logic Refactor
We will refactor `cmd/so-install/main.go` to separate the OS-level `main()` (which calls `os.Exit`) from the application logic `Run()`. This is a standard pattern in Go for testing CLI applications.

```go
func Run(args []string, out io.Writer, errOut io.Writer) int {
    // ... logic ...
}

func main() {
    os.Exit(Run(os.Args, os.Stdout, os.Stderr))
}
```

### Decision: Mocking for Hardware/OS Detection
We will use interfaces for OS detection and shell execution (already in place) to simulate different environments.
For `detectDesktopEnvironment`, we will mock environment variables and the `isProcessRunning` function.

### Decision: TUI Testing
We will use `tea.Model.Update` directly with mock messages to simulate user interaction and verify state changes. For `View`, we will call it and ensure it returns a non-empty string for each state, possibly checking for specific keywords.

## Implementation Plan

### Infrastructure (Nvidia, Desktop, OSRelease)
- Expand existing test suites to inject more diverse mock outputs.
- Specifically, mock `executor.Execute` to return success/failure for specific commands like `ubuntu-drivers` or `apt-cache`.

### Presentation (TUI)
- Add a helper in `model_test.go` to assert state transitions.
- Ensure `Update` is called with `tea.KeyMsg` for each possible interaction.

### Verification Logic
- Use `go test -cover ./...` as the primary verification tool.
- A final task will be to generate a coverage report and confirm the 90% threshold for each package.
