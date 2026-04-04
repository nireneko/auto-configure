# Delta for Infrastructure: ShellExecutor Timeout & Process Group

## MODIFIED Requirements

### Requirement: ShellExecutor — Command Execution with Timeout

The `ShellExecutor` MUST execute commands with a 10-minute `context.WithTimeout`. If the command does not complete within the timeout, the executor MUST kill the process and return an error.

#### Scenario: Command completes within timeout
- GIVEN a shell command that exits normally
- WHEN `Execute()` is called
- THEN it MUST return stdout, stderr, and a nil error

#### Scenario: Command exceeds timeout
- GIVEN a shell command that runs longer than the configured timeout
- WHEN `Execute()` is called
- THEN it MUST return a non-nil error
- AND the process MUST be terminated

### Requirement: ShellExecutor — Process Group Isolation

The `ShellExecutor` MUST set `Setpgid: true` on `SysProcAttr` so that the child process and any processes it spawns run in their own process group.

#### Scenario: Child spawns a daemon process
- GIVEN a shell script that forks a background daemon inheriting stdout/stderr pipes
- WHEN `Execute()` is called
- THEN `cmd.Run()` MUST NOT block indefinitely after the script exits
- AND the executor MUST return once the parent script process exits or the timeout fires

#### Scenario: Daemon does not prevent executor from returning
- GIVEN the Ollama or OpenCode install script starts a background service
- WHEN `Execute()` is called
- THEN the executor MUST return a result (success or timeout error) without hanging

### Requirement: ShellExecutor — Stdout and Stderr Capture Preserved

The executor MUST continue to return trimmed stdout and stderr strings on successful execution. This behavior is UNCHANGED.

#### Scenario: Capture stdout and stderr on success
- GIVEN a command that writes to stdout and stderr
- WHEN `Execute()` completes successfully
- THEN stdout MUST contain the trimmed standard output
- AND stderr MUST contain the trimmed standard error

#### Scenario: Capture stderr on failure
- GIVEN a command that exits with a non-zero code
- WHEN `Execute()` returns an error
- THEN stderr MUST be non-empty and contain the error output
