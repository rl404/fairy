package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
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

// Log is logging client.
type Log struct {
	log   *slog.Logger
	level LogLevel
	json  bool
	color bool
}

// New to create new logging client.
// Color is not working in json format.
func New(level LogLevel, jsonFmt, color bool) *Log {
	log := Log{
		level: level,
		json:  jsonFmt,
		color: color,
	}

	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key != slog.LevelKey {
			return a
		}

		level, ok := a.Value.Any().(slog.Level)
		if !ok {
			return a
		}

		levelLabel := log.getLevelName(level)

		a.Value = slog.StringValue(log.colorize(level, levelLabel))
		return a
	}

	log.log = slog.New(tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
		Level:       convertLevel(level),
		ReplaceAttr: replaceAttr,
		NoColor:     !color,
		TimeFormat:  time.RFC3339,
	}))

	if jsonFmt {
		log.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       convertLevel(level),
			ReplaceAttr: replaceAttr,
		}))
	}

	return &log
}

// Trace to print trace log.
func (l *Log) Trace(str string, args ...interface{}) {
	l.log.Log(context.Background(), levelTrace, fmt.Sprintf(str, args...))
}

// Debug to print debug log.
func (l *Log) Debug(str string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(str, args...))
}

// Info to print info log.
func (l *Log) Info(str string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(str, args...))
}

// Warn to print warn log.
func (l *Log) Warn(str string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(str, args...))
}

// Error to print error log.
func (l *Log) Error(str string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(str, args...))
}

// Fatal to print fatal log.
// Will exit the program when called.
func (l *Log) Fatal(str string, args ...interface{}) {
	l.log.Log(context.Background(), levelFatal, fmt.Sprintf(str, args...))
	os.Exit(1)
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (l *Log) Panic(str string, args ...interface{}) {
	l.log.Log(context.Background(), levelPanic, fmt.Sprintf(str, args...))
	panic(fmt.Sprintf(str, args...))
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (l *Log) Log(fields map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	if level, ok := fields["level"]; ok {
		switch reflect.TypeOf(level).Kind() {
		case reflect.Int8:
			delete(fields, "level")

			attrs := make([]any, 0)
			for k, v := range fields {
				attrs = append(attrs, slog.Any(k, v))
			}

			switch LogLevel(reflect.ValueOf(level).Int()) {
			case TraceLevel:
				l.log.Log(context.Background(), levelTrace, "", attrs...)
			case DebugLevel:
				l.log.Log(context.Background(), slog.LevelDebug, "", attrs...)
			case InfoLevel:
				l.log.Log(context.Background(), slog.LevelInfo, "", attrs...)
			case WarnLevel:
				l.log.Log(context.Background(), slog.LevelWarn, "", attrs...)
			case ErrorLevel:
				l.log.Log(context.Background(), slog.LevelError, "", attrs...)
			case FatalLevel:
				l.log.Log(context.Background(), levelFatal, "", attrs...)
			case PanicLevel:
				l.log.Log(context.Background(), levelPanic, "", attrs...)
			}
		}
	}
}
