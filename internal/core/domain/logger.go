package domain

// Logger defines a simple logging interface to decouple implementation from business logic.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

// NoopLogger is a logger that does nothing. Useful for tests.
type NoopLogger struct{}

func (n NoopLogger) Info(msg string, args ...any)  {}
func (n NoopLogger) Error(msg string, args ...any) {}
func (n NoopLogger) Debug(msg string, args ...any) {}

var _ Logger = (*NoopLogger)(nil)
