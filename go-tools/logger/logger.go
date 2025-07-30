package logger

import (
	"context"
	"log/slog"
	"os"
)

type loggerKeyType string

const (
	loggerCtxKey loggerKeyType = "logger"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(level slog.Level) *Logger {
	return &Logger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})),
	}
}

func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		return logger
	}
	return NewLogger(slog.LevelInfo)
}

func (l *Logger) WarnWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op), slog.String("msg", msg)}, args...)
	l.Warn(msg, args...)
}

func (l *Logger) InfoWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op), slog.String("msg", msg)}, args...)
	l.Info(msg, args...)
}

func (l *Logger) ErrorWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op), slog.String("msg", msg)}, args...)
	l.Error(msg, args...)
}

func (l *Logger) DebugWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op), slog.String("msg", msg)}, args...)
	l.Debug(msg, args...)
}
