package slog

import "log/slog"

const (
	levelTrace = slog.Level(-8)
	levelFatal = slog.Level(9)
	levelPanic = slog.Level(10)
)

var levelNames = map[slog.Level]string{
	levelTrace:      "TRC",
	slog.LevelDebug: "DBG",
	slog.LevelInfo:  "INF",
	slog.LevelWarn:  "WRN",
	slog.LevelError: "ERR",
	levelFatal:      "FTL",
	levelPanic:      "PNC",
}

var levelFullNames = map[slog.Level]string{
	levelTrace:      "TRACE",
	slog.LevelDebug: "DEBUG",
	slog.LevelInfo:  "INFO",
	slog.LevelWarn:  "WARN",
	slog.LevelError: "ERROR",
	levelFatal:      "FATAL",
	levelPanic:      "PANIC",
}

func (l *Log) getLevelName(level slog.Level) string {
	if l.json {
		label, ok := levelFullNames[level]
		if ok {
			return label
		}
		return level.String()
	}

	label, ok := levelNames[level]
	if ok {
		return label
	}
	return level.String()
}

func convertLevel(lvl LogLevel) slog.Level {
	switch lvl {
	case TraceLevel:
		return levelTrace
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case FatalLevel:
		return levelFatal
	case PanicLevel:
		return levelPanic
	default:
		return 11
	}
}
