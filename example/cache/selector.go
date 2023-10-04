// This file contains function to help you switch caching type easily.
//
// You can use this in your project if you want.
//
// You don't need to handle all caching types, just handle
// the ones you will use.
package main

import (
	"errors"
	"time"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/cache/inmemory"
	"github.com/rl404/fairy/cache/memcache"
	"github.com/rl404/fairy/cache/nocache"
	"github.com/rl404/fairy/cache/redis"
)

// CacheType is type for cache.
type CacheType int8

// Available types for cache.
const (
	NoCache CacheType = iota
	InMemory
	Redis
	Memcache
)

// ErrInvalidCacheType is error for invalid cache type.
var ErrInvalidCacheType = errors.New("invalid cache type")

// New to create new cache client depends on the type.
func New(cacheType CacheType, address string, password string, expiredTime time.Duration) (cache.Cacher, error) {
	switch cacheType {
	case NoCache:
		return nocache.New()
	case InMemory:
		return inmemory.New(expiredTime)
	case Redis:
		return redis.New(address, password, expiredTime)
	case Memcache:
		return memcache.New(address, expiredTime)
	default:
		return nil, ErrInvalidCacheType
	}
}
