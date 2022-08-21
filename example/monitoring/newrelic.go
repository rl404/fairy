package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
	_cache "github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/monitoring/newrelic/cache"
	"github.com/rl404/fairy/monitoring/newrelic/database"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Init your newrelic application.
// For best practice, don't declare in global
// variable like this. Instead, pass it in
// param function.
var nrApp *newrelic.Application

func init() {
	nrApp, _ = newrelic.NewApplication(
		newrelic.ConfigAppName("your-app-name"),
		newrelic.ConfigLicense("your-newrelic-license-key"),
		newrelic.ConfigDistributedTracerEnabled(true), // enable tracing
	)
}

func cacheWithNewrelic() {
	ctx := context.Background()

	// Init cache.
	dialect := _cache.Redis
	address := "localhost:6379"
	cacher, err := _cache.New(dialect, address, "", time.Minute)
	if err != nil {
		// Handle error.
		panic(err)
	}

	// Wrap cache with newrelic.
	cacher = cache.New("redis", address, cacher)

	// Start newrelic transaction if not started.
	// You need to do this only once per workflow.
	// For more info, read https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-transactions.
	tx := nrApp.StartTransaction("transaction-name")
	defer tx.End()

	// Put the transaction in context.
	ctx = newrelic.NewContext(ctx, tx)

	// Use cache as usual with context contain newrelic
	// transaction.
	cacher.Set(ctx, "key", "data")
}

func dbWithNewrelic() {
	ctx := context.Background()

	// Init gorm db as usual.
	address := "127.0.0.1:3306"
	dsn := "user:pass@tcp(" + address + ")/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		// Handle error.
		panic(err)
	}

	// Use newrelic as gorm plugin.
	db.Use(database.NewGORM(address, "dbname"))

	// Start newrelic transaction if not started.
	// You need to do this only once per workflow.
	// For more info, read https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-transactions.
	tx := nrApp.StartTransaction("transaction-name")
	defer tx.End()

	// Put the transaction in context.
	ctx = newrelic.NewContext(ctx, tx)

	// Do query as usual with context contain newrelic
	// transaction.
	var model dummyModel
	db.WithContext(ctx).Where("id = ?", 1).First(&model)
}

func httpWithNewrelic() {
	// Init go-chi.
	r := chi.NewRouter()

	// Use newrelic middleware.
	// Automatically start newrelic transaction
	// and put it in the context.
	r.Use(middleware.NewHTTP(nrApp))

	// Register your route.
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		// Get newrelic transaction from context
		// and start segment.
		defer newrelic.FromContext(r.Context()).StartSegment("segment-name").End()

		// Do something.
	})

	// Serve http as usual.
	http.ListenAndServe(":3000", r)
}
