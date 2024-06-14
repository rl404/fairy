// Package redis is a wrapper of the original "github.com/redis/go-redis" library.
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLimit struct {
	client   *redis.Client
	key      string
	rate     int
	interval time.Duration
}

// New to create new redis based limiter.
func New(rate int, interval time.Duration, addressPasswordKey ...string) (*redisLimit, error) {
	address, pass, key := "localhost:6379", "", "fairy:rate-limit"
	for i, v := range addressPasswordKey {
		switch i {
		case 0:
			address = v
		case 1:
			pass = v
		case 2:
			if v != "" {
				key += ":" + v
			}
		}
	}

	// Init redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping test.
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &redisLimit{
		client:   client,
		key:      key,
		rate:     rate,
		interval: interval,
	}, nil
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (r *redisLimit) Take() {
	ctx := context.Background()
	for {
		start := time.Now().UnixNano()
		interval := time.Duration(r.interval.Seconds() * float64(time.Second))
		end := start + int64(interval)

		pipeline := r.client.TxPipeline()
		pipeline.ZRemRangeByScore(ctx, r.key, "-inf", fmt.Sprintf("%d", start))
		pipeline.ZCard(ctx, r.key)

		cmds, err := pipeline.Exec(ctx)
		if err != nil {
			log.Printf("Redis pipeline error: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		currentCount := cmds[1].(*redis.IntCmd).Val()

		if currentCount >= int64(r.rate) {
			waitTime := time.Duration((end-start)/int64(r.rate) - currentCount)
			time.Sleep(waitTime)
			continue
		}

		pipeline.ZAdd(ctx, r.key, redis.Z{Score: float64(end), Member: end})
		pipeline.Exec(ctx)
		break
	}
}
