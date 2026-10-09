// Package httpserver contains HTTP server handlers and middleware.
package httpserver

import (
	"context"
	"io"
	"net/http"
)

const (
	livenessPath  = "/healthz"
	readinessPath = "/readyz"
)

// ReadinessCheck reports whether the application can safely receive traffic.
// Implementations must honor the request context and avoid exposing internal
// dependency errors in the HTTP response.
type ReadinessCheck func(context.Context) error

// RegisterHealthEndpoints registers liveness and readiness endpoints on mux.
// A nil readiness check is treated as not ready so an unconfigured process
// cannot report readiness accidentally.
func RegisterHealthEndpoints(mux *http.ServeMux, check ReadinessCheck) {
	mux.HandleFunc("GET "+livenessPath, func(w http.ResponseWriter, _ *http.Request) {
		writeHealthResponse(w, http.StatusOK, "ok")
	})

	mux.HandleFunc("GET "+readinessPath, func(w http.ResponseWriter, request *http.Request) {
		if check == nil || check(request.Context()) != nil {
			writeHealthResponse(w, http.StatusServiceUnavailable, "not ready")
			return
		}
		writeHealthResponse(w, http.StatusOK, "ready")
	})
}

func writeHealthResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body+"\n")
}
