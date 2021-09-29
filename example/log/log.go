package main

import "github.com/rl404/fairy"

func main() {
	// Log type and level.
	t := fairy.Zerolog
	l := fairy.TraceLevel

	// Init logger.
	log, err := fairy.NewLog(t, l, false, true)
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
