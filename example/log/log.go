package main

import "github.com/rl404/fairy"

func main() {
	// Log type, level, json format, and color.
	t := fairy.Zerolog
	lvl := fairy.TraceLevel
	json := false
	color := true

	// Init logger.
	log, err := fairy.NewLog(t, lvl, json, color)
	if err != nil {
		panic(err)
	}

	log.Trace("%s", "trace")
	log.Debug("%s", "debug")
	log.Info("%s", "info")
	log.Warn("%s", "warn")
	log.Error("%s", "error")
	log.Fatal("%s", "fatal")
	log.Panic("%s", "panic")
}
