package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/rl404/fairy"
)

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

	// General log with additional fields.
	// Key `level` can be used to differentiate
	// log level.
	log.Log(map[string]interface{}{
		"level":  fairy.ErrorLevel,
		"field1": "f1",
		"field2": "f2",
	})

	// Quick log.
	log.Trace("%s", "trace")
	log.Debug("%s", "debug")
	log.Info("%s", "info")
	log.Warn("%s", "warn")
	log.Error("%s", "error")
	// log.Fatal("%s", "fatal")
	// log.Panic("%s", "panic")

	// Example use of go-chi http middleware with log.
	r := chi.NewRouter()
	r.Use(fairy.MiddlewareWithLog(log))

	// Or wrap the handler directly.
	r.Get("/user", fairy.HandlerFuncWithLog(log, sampleHandler))

	// Let's see the printed log.
	// Run this whole main() function to
	// see the log.
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("sample-body"))
	r.Post("/test", sampleHandler)
	r.ServeHTTP(httptest.NewRecorder(), req)
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"ok":"nice"}`)
}
