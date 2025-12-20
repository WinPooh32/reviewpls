package slogutil_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/WinPooh32/reviewpls/internal/slogutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	t.Parallel()

	// Create a test logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a context with the logger
	ctx := slogutil.WithContext(context.Background(), logger)

	// Verify the logger was stored in context by retrieving it
	retrievedLogger := slogutil.Ctx(ctx)
	assert.Same(t, logger, retrievedLogger)
}

func TestCtx(t *testing.T) {
	t.Parallel()

	// Test with a context that has a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := slogutil.WithContext(context.Background(), logger)

	retrievedLogger := slogutil.Ctx(ctx)
	assert.Equal(t, logger, retrievedLogger)

	// Test with a context that doesn't have a logger (should return default)
	emptyCtx := context.Background()

	defaultLogger := slogutil.Ctx(emptyCtx)
	require.NotNil(t, defaultLogger)

	assert.Same(t, slog.Default(), defaultLogger)
}

func TestCtxWithNilLogger(t *testing.T) {
	t.Parallel()

	// Test with a context that has a nil logger
	// We can't directly test this since contextKey is unexported
	// But we can test that it returns the default logger when no valid logger is found
	ctx := context.WithValue(context.Background(), struct{ name string }{name: "logger"}, nil)

	retrievedLogger := slogutil.Ctx(ctx)
	require.NotNil(t, retrievedLogger)

	assert.Same(t, slog.Default(), retrievedLogger)
}
