// Package cache is a prometheus wrapper for "github.com/rl404/fairy/cache"
// cacher interface.
package cache

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rl404/fairy/cache"
)

const (
	cacheReqName     = "cache_requests_total"
	cacheLatencyName = "cache_request_duration_seconds"
	cacheHit         = "HIT"
	cacheMiss        = "MISS"
	cacheGet         = "GET"
	cacheSet         = "SET"
	cacheDelete      = "DELETE"
)

var cp cachePrometheus

type cachePrometheus struct {
	req *prometheus.CounterVec
	lat *prometheus.HistogramVec
}

type client struct {
	dialect string
	cacher  cache.Cacher
}

func init() {
	cp = cachePrometheus{
		req: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: cacheReqName,
				Help: "How many cache requests processed, partitioned by dialect, operation, result and key.",
			},
			[]string{"dialect", "operation", "result", "key"},
		),
		lat: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: cacheLatencyName,
				Help: "How long it took to process the request, partitioned by dialect, operation, result and key.",
			},
			[]string{"dialect", "operation", "result", "key"},
		),
	}

	prometheus.MustRegister(cp.req)
	prometheus.MustRegister(cp.lat)
}

// New to create new prometheus plugin for cache.
func New(d string, c cache.Cacher) cache.Cacher {
	return &client{
		dialect: d,
		cacher:  c,
	}
}

// Get to update get metrics.
func (c *client) Get(key string, data interface{}) error {
	start := time.Now()
	if err := c.cacher.Get(key, data); err != nil {
		cp.req.WithLabelValues(c.dialect, cacheGet, cacheMiss, key).Inc()
		cp.lat.WithLabelValues(c.dialect, cacheGet, cacheMiss, key).Observe(float64(time.Since(start).Seconds()))
		return err
	}
	cp.req.WithLabelValues(c.dialect, cacheGet, cacheHit, key).Inc()
	cp.lat.WithLabelValues(c.dialect, cacheGet, cacheHit, key).Observe(float64(time.Since(start).Seconds()))
	return nil
}

// Set to update set metrics.
func (c *client) Set(key string, data interface{}) error {
	start := time.Now()
	if err := c.cacher.Set(key, data); err != nil {
		cp.req.WithLabelValues(c.dialect, cacheSet, cacheMiss, key).Inc()
		cp.lat.WithLabelValues(c.dialect, cacheSet, cacheMiss, key).Observe(float64(time.Since(start).Seconds()))
		return err
	}
	cp.req.WithLabelValues(c.dialect, cacheSet, cacheHit, key).Inc()
	cp.lat.WithLabelValues(c.dialect, cacheSet, cacheHit, key).Observe(float64(time.Since(start).Seconds()))
	return nil
}

// Delete to update delete metrics.
func (c *client) Delete(key string) error {
	start := time.Now()
	if err := c.cacher.Delete(key); err != nil {
		cp.req.WithLabelValues(c.dialect, cacheDelete, cacheMiss, key).Inc()
		cp.lat.WithLabelValues(c.dialect, cacheDelete, cacheMiss, key).Observe(float64(time.Since(start).Seconds()))
		return err
	}
	cp.req.WithLabelValues(c.dialect, cacheDelete, cacheHit, key).Inc()
	cp.lat.WithLabelValues(c.dialect, cacheDelete, cacheHit, key).Observe(float64(time.Since(start).Seconds()))
	return nil
}

// Close to close.
func (c *client) Close() error {
	return c.cacher.Close()
}