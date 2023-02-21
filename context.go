package kvlog

import "context"

type contextKeyType string

const contextKey contextKeyType = "kvlog.logger"

// ContextWithLogger creates a new context.Context derived from ctx that
// contains l as a value.
// The logger can be retrieved later by calling FromContext.
func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, contextKey, l)
}

// FromContext returns the Logger associated with ctx. If no logger is
// associated with ctx, a NoOpLogger is returned.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(contextKey).(Logger); ok {
		return l
	}

	return noOpLoggerValue
}
