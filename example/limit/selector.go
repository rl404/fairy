// This file contains function to help you switch limiter type easily.
//
// You can use this in your project if you want.
//
// You don't need to handle all limiter types, just handle
// the ones you will use.
package main

import (
	"errors"
	"time"

	"github.com/rl404/fairy/limit"
	"github.com/rl404/fairy/limit/atomic"
	"github.com/rl404/fairy/limit/builtin"
	"github.com/rl404/fairy/limit/mutex"
	"github.com/rl404/fairy/limit/redis"
)

// LimitType is type for rate limit.
type LimitType int8

// Available types for rate limit.
const (
	Mutex LimitType = iota
	Atomic
	Builtin
	Redis
)

// ErrInvalidLimitType is error for invalid rate limit type.
var ErrInvalidLimitype = errors.New("invalid rate limit type")

// New to create new rate limiter.
func New(limitType LimitType, rate int, interval time.Duration, options ...string) (limit.Limiter, error) {
	switch limitType {
	case Mutex:
		return mutex.New(rate, interval), nil
	case Atomic:
		return atomic.New(rate, interval), nil
	case Builtin:
		return builtin.New(rate, interval), nil
	case Redis:
		return redis.New(rate, interval, options...)
	default:
		return nil, ErrInvalidLimitype
	}
}
