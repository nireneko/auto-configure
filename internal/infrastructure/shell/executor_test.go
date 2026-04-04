package shell_test

import (
	"testing"
	"time"

	"github.com/so-install/internal/infrastructure/shell"
)

func TestShellExecutor_SuccessfulCommand(t *testing.T) {
	ex := shell.NewShellExecutor()
	stdout, stderr, err := ex.Execute("echo", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stdout != "hello world" {
		t.Errorf("expected 'hello world', got %q", stdout)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
}

func TestShellExecutor_FailingCommand(t *testing.T) {
	ex := shell.NewShellExecutor()
	_, _, err := ex.Execute("false")
	if err == nil {
		t.Fatal("expected error from 'false' command, got nil")
	}
}

func TestShellExecutor_CapturesStderr(t *testing.T) {
	ex := shell.NewShellExecutor()
	// sh -c "echo errtext >&2; exit 1" writes to stderr and exits non-zero
	_, stderr, err := ex.Execute("sh", "-c", "echo errtext >&2; exit 1")
	if err == nil {
		t.Fatal("expected error")
	}
	if stderr != "errtext" {
		t.Errorf("expected 'errtext' in stderr, got %q", stderr)
	}
}

func TestShellExecutor_Timeout(t *testing.T) {
	ex := shell.NewShellExecutorWithTimeout(500 * time.Millisecond)
	start := time.Now()
	_, _, err := ex.Execute("sh", "-c", "sleep 999")
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if elapsed > 3*time.Second {
		t.Errorf("Execute took too long (%v), expected ~500ms", elapsed)
	}
}

func TestShellExecutor_DaemonDoesNotHang(t *testing.T) {
	// sh exits immediately after forking sleep to the background.
	// sleep inherits the stdout/stderr pipes, which would block cmd.Run() forever
	// without WaitDelay. With WaitDelay = 500ms, Execute returns after ~500ms.
	ex := shell.NewShellExecutorWithTimeout(500 * time.Millisecond)
	start := time.Now()
	_, _, _ = ex.Execute("sh", "-c", "sleep 999 &")
	elapsed := time.Since(start)
	if elapsed > 3*time.Second {
		t.Errorf("Execute blocked for %v — daemon kept the pipe open past WaitDelay", elapsed)
	}
}
