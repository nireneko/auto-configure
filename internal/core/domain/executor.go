package domain

// Executor abstracts shell command execution for testability.
type Executor interface {
	Execute(name string, args ...string) (stdout, stderr string, err error)
}
