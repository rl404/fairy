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
)

// Example use of go-chi http middleware with log.
func httpWithLog(l log.Logger) {
	r := chi.NewRouter()
	r.Use(log.HTTPMiddlewareWithLog(l, log.APIMiddlewareConfig{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
		RawPath:        true,
		Error:          true,
	}))

	// Or wrap the handler directly.
	r.Get("/user", log.HTTPHandlerFuncWithLog(l, sampleHandler))

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
