package shell_test

import (
	"testing"

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
