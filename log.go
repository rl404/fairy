package fairy

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
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
)

// ErrInvalidLogType is error for invalid log type.
var ErrInvalidLogType = errors.New("invalid log type")

// NewLog to create new log client depends on the type.
// Color is not working in json format.
func NewLog(logType LogType, level LogLevel, jsonFormat bool, color bool) (Logger, error) {
	switch logType {
	case NoLog:
		return nil, nil
	case BuiltIn:
		return nil, nil
	case Zerolog:
		return zerolog.New(zerolog.LogLevel(level), jsonFormat, color), nil
	case Logrus:
		return nil, nil
	default:
		return nil, ErrInvalidLogType
	}
}

// MiddlewareWithLog is http middleware that will log the request and response.
func MiddlewareWithLog(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return HandlerWithLog(logger, next)
	}
}

// HandlerFuncWithLog is http handler func with log.
func HandlerFuncWithLog(logger Logger, next http.HandlerFunc) http.HandlerFunc {
	return HandlerWithLog(logger, next).(http.HandlerFunc)
}

// HandlerWithLog is http handler with log.
// Also includes error stack tracing feature
// if you use it.
func HandlerWithLog(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Prepare error stack tracing.
		s := NewErrStacker()
		ctx := s.Init(r.Context())

		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r.WithContext(ctx))

		m := map[string]interface{}{
			"level":    getLevelFromStatus(ww.Status()),
			"duration": time.Since(start).String(),
			"method":   r.Method,
			"path":     r.RequestURI,
			"code":     ww.Status(),
			"ip":       getIP(r),
		}

		// Include the error stack if you use it.
		errStack := s.Get(ctx).([]string)
		if len(errStack) > 0 {
			// Reverse the stack order.
			for i, j := 0, len(errStack)-1; i < j; i, j = i+1, j-1 {
				errStack[i], errStack[j] = errStack[j], errStack[i]
			}
			m["error"] = strings.Join(errStack, " | ")
		}

		logger.Log(m)
	})
}

func getLevelFromStatus(status int) LogLevel {
	switch status {
	case
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		return InfoLevel
	case
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusPreconditionFailed,
		http.StatusExpectationFailed,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusFailedDependency,
		http.StatusTooManyRequests:
		return WarnLevel
	default:
		return ErrorLevel
	}
}

func getIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
