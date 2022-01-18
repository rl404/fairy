package log

type loggerChain struct {
	loggers []Logger
}

// NewChain to create new logger chain.
// Useful if you want to print the log in
// local and send the log to third party at
// the same time.
func NewChain(logger Logger, loggers ...Logger) Logger {
	return &loggerChain{
		loggers: append([]Logger{logger}, loggers...),
	}
}

// Trace to print trace log.
func (lc *loggerChain) Trace(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Trace(format, args...)
	}
}

// Debug to print debug log.
func (lc *loggerChain) Debug(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Debug(format, args...)
	}
}

// Info to print info log.
func (lc *loggerChain) Info(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Info(format, args...)
	}
}

// Warn to print warn log.
func (lc *loggerChain) Warn(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Warn(format, args...)
	}
}

// Error to print error log.
func (lc *loggerChain) Error(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Error(format, args...)
	}
}

// Fatal to print fatal log.
func (lc *loggerChain) Fatal(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Fatal(format, args...)
	}
}

// Panic to print panic log.
func (lc *loggerChain) Panic(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Panic(format, args...)
	}
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (lc *loggerChain) Log(fields map[string]interface{}) {
	for _, l := range lc.loggers {
		l.Log(fields)
	}
}
