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

func NewLogger(env string) *Logger {
	var handler slog.Handler

	switch env {
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default: // local
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		return logger
	}

	if env, ok := ctx.Value("ENV").(string); ok {
		return NewLogger(env)
	}

	if env, ok := os.LookupEnv("ENV"); ok {
		return NewLogger(env)
	}

	return NewLogger("local")
}

func (l *Logger) WarnWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op)}, args...)
	l.Warn(msg, args...)
}

func (l *Logger) InfoWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op)}, args...)
	l.Info(msg, args...)
}

func (l *Logger) ErrorWithOp(op, msg string, err error, args ...any) {
	args = append([]any{slog.String("op", op), slog.Any("err", err)}, args...)
	l.Error(msg, args...)
}

func (l *Logger) DebugWithOp(op, msg string, args ...any) {
	args = append([]any{slog.String("op", op)}, args...)
	l.Debug(msg, args...)
}
