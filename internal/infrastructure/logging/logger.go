package logging

import (
	"context"
	"io"
	"log/slog"

	"github.com/so-install/internal/core/domain"
)

// FileLogger implements domain.Logger using slog and a file handler.
type FileLogger struct {
	logger *slog.Logger
}

// NewFileLogger creates a new FileLogger that writes to the given io.Writer.
func NewFileLogger(w io.Writer) *FileLogger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &FileLogger{
		logger: slog.New(handler),
	}
}

var _ domain.Logger = (*FileLogger)(nil)

func (f *FileLogger) Info(msg string, args ...any) {
	f.logger.Info(msg, args...)
}

func (f *FileLogger) Error(msg string, args ...any) {
	f.logger.Error(msg, args...)
}

func (f *FileLogger) Debug(msg string, args ...any) {
	f.logger.Debug(msg, args...)
}

// LogWriter wraps the slog logger to implement io.Writer,
// for use with exec.Cmd's Stdout/Stderr if needed.
type LogWriter struct {
	logger *slog.Logger
	level  slog.Level
}

func NewLogWriter(logger *slog.Logger, level slog.Level) *LogWriter {
	return &LogWriter{logger: logger, level: level}
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.logger.Log(context.Background(), w.level, string(p))
	return len(p), nil
}
