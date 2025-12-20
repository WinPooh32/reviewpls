package slogutil

import (
	"context"
	"log/slog"
)

// contextKey is the key used to store and retrieve the logger from context.
type contextKey struct{}

// WithContext stores the logger in the context
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// Ctx retrieves the logger from the context
func Ctx(ctx context.Context) *slog.Logger {
	logger := ctx.Value(contextKey{})
	if logger != nil {
		if l, ok := logger.(*slog.Logger); ok {
			return l
		}
	}
	// Return a default logger if none is found in context
	return slog.Default()
}
