// Package builtin is a wrapper of the built-in "golang.org/x/time/rate" library.
package builtin

import (
	"context"
	"time"

	_rate "golang.org/x/time/rate"
)

type rateLimit struct {
	limiter *_rate.Limiter
}

// New to create new built-in limiter.
func New(rate int, interval time.Duration) *rateLimit {
	ratePerSecond := float64(rate) / interval.Seconds()
	return &rateLimit{
		limiter: _rate.NewLimiter(_rate.Limit(ratePerSecond), rate),
	}
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (r *rateLimit) Take() {
	r.limiter.Wait(context.Background())
}
