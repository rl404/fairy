package log

import (
	"context"
	"time"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/fairy/pubsub"
)

// PubSubMiddlewareConfig is pubsub config for middleware.
type PubSubMiddlewareConfig struct {
	// Topic/queue name.
	Topic string
	// Show message payload.
	Payload bool
	// Show error stack.
	Error bool
}

// PubSubMiddlewareWithLog is pubsub middelware that will log the request.
func PubSubMiddlewareWithLog(logger Logger, middlewareConfig ...PubSubMiddlewareConfig) func(pubsub.HandlerFunc) pubsub.HandlerFunc {
	return func(next pubsub.HandlerFunc) pubsub.HandlerFunc {
		return PubSubHandlerFuncWithLog(logger, next, middlewareConfig...)
	}
}

// PubSubHandlerFuncWithLog is pubsub handler with log.
// Also includes error stack tracing feature
// if you use it.
func PubSubHandlerFuncWithLog(logger Logger, next pubsub.HandlerFunc, middlewareConfig ...PubSubMiddlewareConfig) pubsub.HandlerFunc {
	var cfg PubSubMiddlewareConfig
	if len(middlewareConfig) > 0 {
		cfg = middlewareConfig[0]
	}

	return func(ctx context.Context, message []byte) {
		if logger == nil {
			next(ctx, message)
			return
		}

		// Prepare error stack tracing.
		ctx = errors.Init(ctx)
		start := time.Now()

		// Call next handler.
		next(ctx, message)

		// Prepare map for logging.
		m := map[string]interface{}{
			"level":    infoLevel,
			"duration": time.Since(start).String(),
		}

		if cfg.Topic != "" {
			m["topic"] = cfg.Topic
		}

		if cfg.Payload {
			m["payload"] = string(message)
		}

		// Include the error stack if you use it.
		errStack := errors.Get(ctx)
		if len(errStack) > 0 {
			m["level"] = errorLevel

			if cfg.Error {
				// Copy slice to prevent reversed multiple times
				// if using multiple middleware.
				errTmp := cpSlice(errStack)

				// Reverse the stack order.
				for i, j := 0, len(errTmp)-1; i < j; i, j = i+1, j-1 {
					errTmp[i], errTmp[j] = errTmp[j], errTmp[i]
				}

				m["error"] = errTmp
			}
		}

		logger.Log(m)
	}
}
