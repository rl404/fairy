// This file contains function to help you switch logging type easily.
//
// You can use this in your project if you want.
//
// You don't need to handle all logging types, just handle
// the ones you will use.
package main

import (
	"errors"

	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/log/elasticsearch"
	"github.com/rl404/fairy/log/logrus"
	"github.com/rl404/fairy/log/newrelic"
	"github.com/rl404/fairy/log/slog"
	"github.com/rl404/fairy/log/zap"
	"github.com/rl404/fairy/log/zerolog"
)

// LogLevel is level of log that will be printed.
// Will print level that is higher than your
// chosen one.
type LogLevel int8

// Available log level.
const (
	TraceLevel LogLevel = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	Disabled
)

// LogType is type for logging library.
type LogType int8

// Available types for logging.
const (
	Zerolog LogType = iota + 1
	Logrus
	Zap
	Elasticsearch
	Newrelic
	Slog
)

// ErrInvalidLogType is error for invalid log type.
var ErrInvalidLogType = errors.New("invalid log type")

// Config is log config.
type Config struct {
	Type       LogType
	Level      LogLevel
	JsonFormat bool
	Color      bool

	// For elasticsearch.
	ElasticsearchAddresses []string
	ElasticsearchUser      string
	ElasticsearchPassword  string
	ElasticsearchIndex     string
	ElasticsearchIsSync    bool

	// For newrelic.
	NewrelicName    string
	NewrelicLicense string
}

// New to create new log client depends on the type.
// Color will not work in json format.
func New(cfg Config) (log.Logger, error) {
	switch cfg.Type {
	case Zerolog:
		return zerolog.New(zerolog.LogLevel(cfg.Level), cfg.JsonFormat, cfg.Color), nil
	case Logrus:
		return logrus.New(logrus.LogLevel(cfg.Level), cfg.JsonFormat, cfg.Color), nil
	case Zap:
		return zap.New(zap.LogLevel(cfg.Level), cfg.JsonFormat, cfg.Color)
	case Elasticsearch:
		return elasticsearch.New(elasticsearch.Config{
			Addresses: cfg.ElasticsearchAddresses,
			Username:  cfg.ElasticsearchUser,
			Password:  cfg.ElasticsearchPassword,
			Index:     cfg.ElasticsearchIndex,
			Level:     elasticsearch.LogLevel(cfg.Level),
			IsSync:    cfg.ElasticsearchIsSync,
		})
	case Newrelic:
		return newrelic.New(newrelic.Config{
			Name:       cfg.NewrelicName,
			LicenseKey: cfg.NewrelicLicense,
			Level:      newrelic.LogLevel(cfg.Level),
		})
	case Slog:
		return slog.New(slog.LogLevel(cfg.Level), cfg.JsonFormat, cfg.Color), nil
	default:
		return nil, ErrInvalidLogType
	}
}
