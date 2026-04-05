# Specification: Fix Shell Executor WaitDelay Hang

## Requirements

### 1. WaitDelay Isolation
The shell executor MUST NOT use the main command timeout for the `WaitDelay`.

### 2. Default WaitDelay
The shell executor SHALL use a default `WaitDelay` of 5 seconds when not specified otherwise.

### 3. Timeout Preservation
The main command timeout (10 minutes by default) MUST remain unchanged and still be used to kill hanging foreground processes.

## Scenarios

### Scenario 1: Successful Short Command
**Given** a command that completes in 100ms
**When** Execute is called
**Then** it MUST return within ~100ms
**And** stdout MUST be captured correctly

### Scenario 2: Foreground Hang (Timeout)
**Given** a command that sleeps for 999s
**And** a timeout of 500ms
**When** Execute is called
**Then** it MUST timeout and kill the process group within ~500ms
**And** an error MUST be returned

### Scenario 3: Daemon Spawned (WaitDelay)
**Given** a command that forks a background process (e.g., `sleep 999 &`) and exits immediately
**And** a WaitDelay of 500ms
**And** a timeout of 10 minutes
**When** Execute is called
**Then** the command process exits immediately
**And** the executor MUST wait at most ~500ms for the pipes to close
**And** it MUST return after ~500ms
**And** no error MUST be returned (as the main process exited with 0)
