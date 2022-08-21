// Package database is a prometheus wrapper for database.
package database

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

const gormReqName = "database_requests_total"
const gormLatencyName = "database_request_duration_seconds"

type gormPlugin struct {
	dbName string
	req    *prometheus.CounterVec
	lat    *prometheus.HistogramVec
}

// NewGORM to create new prometheus plugin for gorm.
func NewGORM(dbName string) *gormPlugin {
	gp := gormPlugin{
		dbName: dbName,
		req: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: gormReqName,
				Help: "How many database queries processed, partitioned by dialect, database, operation, table and query.",
			},
			[]string{"dialect", "database", "operation", "table", "query"},
		),
		lat: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: gormLatencyName,
				Help: "How long it took to process the query, partitioned by dialect, database, operation, table and query.",
			},
			[]string{"dialect", "database", "operation", "table", "query"},
		),
	}

	prometheus.MustRegister(gp.req)
	prometheus.MustRegister(gp.lat)

	return &gp
}

// Name is plugin name.
func (gp *gormPlugin) Name() string {
	return "prometheus"
}

// Initialize to init newrelic plugin.
func (gp *gormPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Before("gorm:query").Register("prometheus:before_query", gp.start("SELECT"))
	db.Callback().Query().After("gorm:query").Register("prometheus:after_query", gp.end("SELECT"))

	db.Callback().Create().Before("gorm:create").Register("prometheus:before_create", gp.start("INSERT"))
	db.Callback().Create().After("gorm:create").Register("prometheus:after_create", gp.end("INSERT"))

	db.Callback().Delete().Before("gorm:delete").Register("prometheus:before_delete", gp.start("DELETE"))
	db.Callback().Delete().After("gorm:delete").Register("prometheus:after_delete", gp.end("DELETE"))

	db.Callback().Update().Before("gorm:update").Register("prometheus:before_update", gp.start("UPDATE"))
	db.Callback().Update().After("gorm:update").Register("prometheus:after_update", gp.end("UPDATE"))

	db.Callback().Row().Before("gorm:row").Register("prometheus:before_row", gp.start("ROW"))
	db.Callback().Row().After("gorm:row").Register("prometheus:after_row", gp.start("ROW"))

	db.Callback().Raw().Before("gorm:raw").Register("prometheus:before_raw", gp.start("RAW"))
	db.Callback().Raw().After("gorm:raw").Register("prometheus:after_raw", gp.end("RAW"))

	return nil
}

func (gp *gormPlugin) key(operation string) string {
	return fmt.Sprintf("prometheus:startKey:%s", operation)
}

func (gp *gormPlugin) start(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		db.Set(gp.key(operation), time.Now())
	}
}

func (gp *gormPlugin) end(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		start, ok := db.Get(gp.key(operation))
		if !ok {
			return
		}

		table := db.Statement.Table
		query := db.Statement.SQL.String()

		gp.req.WithLabelValues(db.Name(), gp.dbName, operation, table, query).Inc()
		gp.lat.WithLabelValues(db.Name(), gp.dbName, operation, table, query).Observe(float64(time.Since(start.(time.Time)).Seconds()))
	}
}
