package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WinPooh32/reviewpls/internal/retry"
	"github.com/stretchr/testify/assert"
)

func TestRun_Success(t *testing.T) {
	t.Parallel()

	var count int

	fn := func() error {
		count++
		return nil
	}

	err := retry.Run(context.Background(), fn)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestRun_FatalError(t *testing.T) {
	t.Parallel()

	fn := func() error {
		return retry.ErrFatal
	}

	err := retry.Run(context.Background(), fn)
	assert.ErrorIs(t, err, retry.ErrFatal)
}

func TestRun_MaxAttemptsReached(t *testing.T) {
	t.Parallel()

	var count int

	fn := func() error {
		count++
		return errors.New("some error")
	}

	err := retry.Run(context.Background(), fn, retry.WithMaxAttempts(3))

	assert.ErrorIs(t, err, retry.ErrMaxAttemptsLimitReached)
	assert.Equal(t, 4, count) // 3 retries + 1 success
}

func TestRun_Delay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var count int

	fn := func() error {
		count++
		return errors.New("some error")
	}

	start := time.Now()
	err := retry.Run(ctx, fn, retry.WithDelay(50*time.Millisecond), retry.WithMaxAttempts(2))
	duration := time.Since(start)

	assert.ErrorIs(t, err, retry.ErrMaxAttemptsLimitReached)
	assert.GreaterOrEqual(t, duration, 100*time.Millisecond)
	assert.Equal(t, 3, count)
}

func TestRun_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fn := func() error {
		return errors.New("some error")
	}

	err := retry.Run(ctx, fn)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestRun_SuccessAfterRetry(t *testing.T) {
	t.Parallel()

	var count int

	fn := func() error {
		count++
		if count < 3 {
			return errors.New("some error")
		}

		return nil
	}

	err := retry.Run(context.Background(), fn, retry.WithDelay(1*time.Millisecond))

	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
