package log

import (
	"errors"

	"github.com/rl404/fairy/log/builtin"
	"github.com/rl404/fairy/log/logrus"
	"github.com/rl404/fairy/log/nolog"
	"github.com/rl404/fairy/log/zap"
	"github.com/rl404/fairy/log/zerolog"
)

// Logger is logging interface.
//
// See usage example in example folder.
type Logger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})

	// General log with key value.
	Log(fields map[string]interface{})
}

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
	NoLog LogType = iota
	BuiltIn
	Zerolog
	Logrus
	Zap
)

// ErrInvalidLogType is error for invalid log type.
var ErrInvalidLogType = errors.New("invalid log type")

// New to create new log client depends on the type.
// Color will not work in json format.
func New(logType LogType, level LogLevel, jsonFormat bool, color bool) (Logger, error) {
	switch logType {
	case NoLog:
		return nolog.New(), nil
	case BuiltIn:
		return builtin.New(builtin.LogLevel(level), jsonFormat, color), nil
	case Zerolog:
		return zerolog.New(zerolog.LogLevel(level), jsonFormat, color), nil
	case Logrus:
		return logrus.New(logrus.LogLevel(level), jsonFormat, color), nil
	case Zap:
		return zap.New(zap.LogLevel(level), jsonFormat, color), nil
	default:
		return nil, ErrInvalidLogType
	}
}
