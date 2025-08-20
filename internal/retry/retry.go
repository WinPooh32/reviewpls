package retry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WinPooh32/opt"
)

var (
	ErrFatal                   = errors.New("fatal")
	ErrMaxAttemptsLimitReached = errors.New("max attempts limit is reached")
)

type options struct {
	delay       opt.T[time.Duration]
	maxAttempts opt.T[int]
}

type OptionFunc func(o *options)

func WithDelay(d time.Duration) OptionFunc {
	return func(o *options) {
		o.delay = opt.Wrap(d)
	}
}

func WithMaxAttempts(n int) OptionFunc {
	return func(o *options) {
		o.maxAttempts = opt.Wrap(n)
	}
}

func Run(ctx context.Context, fn func() error, opts ...OptionFunc) error {
	var o options
	for _, optFn := range opts {
		optFn(&o)
	}

	for i := 0; ; i++ {
		if o.delay.Set() && i > 0 {
			select {
			case <-time.After(o.delay.Value()):
				fmt.Println("delay", o.delay.Value())
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		if o.maxAttempts.Set() && i > o.maxAttempts.Value() {
			return ErrMaxAttemptsLimitReached
		}

		err := fn()
		if errors.Is(err, ErrFatal) {
			return err
		}

		if err != nil {
			continue
		}

		return nil
	}
}
