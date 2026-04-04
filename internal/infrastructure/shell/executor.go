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

const defaultTimeout = 10 * time.Minute

// ShellExecutor runs real shell commands using os/exec.
type ShellExecutor struct {
	timeout time.Duration
}

// NewShellExecutor creates a ShellExecutor with the default 10-minute timeout.
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{timeout: defaultTimeout}
}

// NewShellExecutorWithTimeout creates a ShellExecutor with a custom timeout.
// Intended for use in tests.
func NewShellExecutorWithTimeout(d time.Duration) *ShellExecutor {
	return &ShellExecutor{timeout: d}
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
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = e.timeout
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), err
}
