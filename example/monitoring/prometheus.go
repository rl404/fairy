package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_cache "github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/monitoring/prometheus/cache"
	"github.com/rl404/fairy/monitoring/prometheus/database"
	"github.com/rl404/fairy/monitoring/prometheus/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Don't forget to expose your Prometheus metrics so Prometheus can scrape them.
// https://prometheus.io/docs/guides/go-application/

type dummyModel struct{}

func dbWithPrometheus() {
	// Init gorm db as usual.
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		// Handle error.
	}

	// Register prometheus to db.
	database.RegisterGORM("dbname", db)

	// Do query as usual.
	var model dummyModel
	db.Where("id = ?", 1).First(&model)

	// Sample metrics.
	// database_requests_total{operation="INSERT",table="user"} 68
	// database_requests_total{operation="SELECT",table="user"} 2370
	// database_requests_total{operation="UPDATE",table="user"} 216
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.005"} 5947
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.01"} 6180
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.025"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.05"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.1"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.25"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="0.5"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="1"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="2.5"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="5"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="10"} 6214
	// database_request_duration_seconds_bucket{operation="UPDATE",table="user",le="+Inf"} 6214
	// database_request_duration_seconds_sum{operation="UPDATE",table="user"} 12.041487635000017
	// database_request_duration_seconds_count{operation="UPDATE",table="user"} 6214
}

func cacheWithPrometheus() {
	// Init cache.
	dialect := _cache.Redis
	cacher, err := _cache.New(dialect, "localhost:6379", "", time.Minute)
	if err != nil {
		// Handle error.
	}

	// Wrap cache with prometheus.
	cacher = cache.New("redis", cacher)

	// Use cache as usual.
	cacher.Set(context.Background(), "key", "data")

	// Sample metrics.
	// cache_requests_total{dialect="redis",operation="GET",result="MISS"} 1
	// cache_requests_total{dialect="redis",operation="SET",result="HIT"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.005"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.01"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.025"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.05"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.1"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.25"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="0.5"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="1"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="2.5"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="5"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="10"} 1
	// cache_request_duration_seconds_bucket{dialect="redis",operation="SET",result="HIT",le="+Inf"} 1
	// cache_request_duration_seconds_sum{dialect="redis",operation="SET",result="HIT"} 0.000385649
	// cache_request_duration_seconds_count{dialect="redis",operation="SET",result="HIT"} 1
}

func httpWithPrometheus() {
	// Init go-chi.
	r := chi.NewRouter()

	// Use prometheus middleware.
	r.Use(middleware.NewHTTP())

	// Register metrics endpoint, so prometheus can scrape the metrics.
	r.Handle("/metrics", promhttp.Handler())

	// Serve http as usual.
	http.ListenAndServe(":3000", r)

	// Sample metrics.
	// http_requests_total{code="200",method="GET",path="/ping"} 38
	// http_requests_total{code="200",method="POST",path="/user"} 3
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.005"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.01"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.025"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.05"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.1"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.25"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="0.5"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="1"} 0
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="2.5"} 1
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="5"} 1
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="10"} 1
	// http_request_duration_seconds_bucket{code="500",method="POST",path="/register",le="+Inf"} 1
	// http_request_duration_seconds_sum{code="500",method="POST",path="/register"} 1.603506212
	// http_request_duration_seconds_count{code="500",method="POST",path="/register"} 1
}

func main() {}
