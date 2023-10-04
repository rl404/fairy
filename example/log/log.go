package main

import (
	"bytes"
	"context"
	_errors "errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/log/chain"
)

func main() {
	// Init first logger.
	l1, err := New(Config{
		Type:       Zerolog,
		Level:      TraceLevel,
		JsonFormat: false,
		Color:      true,
	})
	if err != nil {
		panic(err)
	}

	// Init second logger.
	l2, err := New(Config{
		Type:                   Elasticsearch,
		Level:                  ErrorLevel,
		ElasticsearchAddresses: []string{"http://localhost:9200"},
		ElasticsearchUser:      "elastic",
		ElasticsearchPassword:  "",
		ElasticsearchIndex:     "logs-example",
		ElasticsearchIsSync:    true,
	})
	if err != nil {
		panic(err)
	}

	// Chain the loggers.
	l := chain.New(l1, l2)

	// General log with additional fields.
	// Key `level` can be used to differentiate
	// log level.
	l.Log(map[string]interface{}{
		"level":  ErrorLevel,
		"field1": "f1",
		"field2": "f2",
	})

	// Quick log.
	l.Trace("%s", "trace")
	l.Debug("%s", "debug")
	l.Info("%s", "info")
	l.Warn("%s", "warn")
	l.Error("%s", "error")
	// l.Fatal("%s", "fatal")
	// l.Panic("%s", "panic")

	// Example use of go-chi http middleware with log.
	r := chi.NewRouter()
	r.Use(log.MiddlewareWithLog(l, log.MiddlewareConfig{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
		RawPath:        true,
		Error:          true,
	}))

	// Or wrap the handler directly.
	r.Get("/user", log.HandlerFuncWithLog(l, sampleHandler))

	// Let's see the printed log.
	// Run this whole main() function to
	// see the log.
	req := httptest.NewRequest(http.MethodPost, "/test/123?query1=a&query2=b", bytes.NewBufferString(`{"name":"sample-request"}`))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	r.Post("/test/{id}", sampleHandler)
	r.ServeHTTP(httptest.NewRecorder(), req)
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	if err := sampleErr(r.Context()); err != nil {
		// Let's also test the error stack trace feature.
		errors.Wrap(r.Context(), _errors.New("sample error"), err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"ok":"nice"}`)
}

func sampleErr(ctx context.Context) error {
	return errors.Wrap(ctx, _errors.New("sample original error"))
}
