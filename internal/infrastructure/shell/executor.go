package shell

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/so-install/internal/core/domain"
)

// ShellExecutor runs real shell commands using os/exec.
type ShellExecutor struct{}

// NewShellExecutor creates a new ShellExecutor.
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

var _ domain.Executor = (*ShellExecutor)(nil)

// Execute runs the named command with args and returns trimmed stdout, stderr, and error.
func (e *ShellExecutor) Execute(name string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), err
}
