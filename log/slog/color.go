package slog

import (
	"fmt"
	"log/slog"
)

var colors = map[slog.Level]int{
	levelTrace:      94,
	slog.LevelDebug: 93,
	slog.LevelInfo:  32,
	slog.LevelWarn:  33,
	slog.LevelError: 91,
	levelFatal:      31,
	levelPanic:      35,
}

func (l *Log) colorize(level slog.Level, s string) string {
	if !l.color || l.json {
		return s
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", colors[level], s)
}
