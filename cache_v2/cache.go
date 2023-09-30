package cache

import (
	"context"
	"time"
)

// Cacher is caching interface.
//
// See usage example in example folder.
type Cacher interface {
	// Get data from cache. The returned value will be
	// assigned to param `data`. Param `data` should
	// be a pointer just like when using json.Unmarshal.
	Get(ctx context.Context, key string, data interface{}) error
	// Save data to cache. Set and Get should be using
	// the same encoding method. For example, json.Marshal
	// for Set and json.Unmarshal for Get.
	Set(ctx context.Context, key string, data interface{}, ttl ...time.Duration) error
	// Delete data from cache.
	Delete(ctx context.Context, key string) error
	// Close cache connection.
	Close() error
}
