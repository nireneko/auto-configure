package shell

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/so-install/internal/core/domain"
)

const (
	defaultTimeout   = 10 * time.Minute
	defaultWaitDelay = 5 * time.Second
)

// ShellExecutor runs real shell commands using os/exec.
type ShellExecutor struct {
	timeout   time.Duration
	waitDelay time.Duration
	logger    domain.Logger
}

// NewShellExecutor creates a ShellExecutor with the default 10-minute timeout,
// 5-second wait delay, and a logger.
func NewShellExecutor(logger domain.Logger) *ShellExecutor {
	if logger == nil {
		logger = domain.NoopLogger{}
	}
	return &ShellExecutor{
		timeout:   defaultTimeout,
		waitDelay: defaultWaitDelay,
		logger:    logger,
	}
}

// NewShellExecutorWithTimeout creates a ShellExecutor with a custom timeout and wait delay.
// Intended for use in tests.
func NewShellExecutorWithTimeout(timeout, waitDelay time.Duration, logger domain.Logger) *ShellExecutor {
	if logger == nil {
		logger = domain.NoopLogger{}
	}
	return &ShellExecutor{
		timeout:   timeout,
		waitDelay: waitDelay,
		logger:    logger,
	}
}

var _ domain.Executor = (*ShellExecutor)(nil)

// Execute runs the named command with args and returns trimmed stdout, stderr, and error.
//
// Two mechanisms prevent hanging:
//
//  1. context.WithTimeout — if the command does not exit within the timeout, the
//     process group is killed via cmd.Cancel.
//
//  2. cmd.WaitDelay — if the command exits normally but a child daemon has inherited
//     the stdout/stderr pipes (keeping them open), WaitDelay forces the I/O goroutines
//     to close after the same timeout, unblocking cmd.Run().
//
// Setpgid isolates the child in its own process group so that cmd.Cancel can send
// SIGKILL to the whole group (including any daemonized grandchildren).
func (e *ShellExecutor) Execute(name string, args ...string) (stdout, stderr string, err error) {
	e.logger.Info("executing command", "command", name, "args", args)

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = e.waitDelay
	cmd.Cancel = func() error {
		e.logger.Error("command timeout, killing process group", "command", name)
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()

	stdout = strings.TrimSpace(outBuf.String())
	stderr = strings.TrimSpace(errBuf.String())

	if err != nil {
		e.logger.Error("command execution failed", "command", name, "err", err, "stderr", stderr)
	} else {
		e.logger.Debug("command execution successful", "command", name, "stdout", stdout)
	}

	return stdout, stderr, err
}
