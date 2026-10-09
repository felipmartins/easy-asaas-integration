package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLivenessDoesNotDependOnReadiness(t *testing.T) {
	checkCalled := false
	handler := healthMux(func(context.Context) error {
		checkCalled = true
		return errors.New("database unavailable")
	})

	response := makeRequest(handler, http.MethodGet, livenessPath, context.Background())

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "ok\n" {
		t.Errorf("body = %q, want %q", response.Body.String(), "ok\n")
	}
	if checkCalled {
		t.Fatal("liveness endpoint called the readiness check")
	}
	assertHealthHeaders(t, response)
}

func TestReadinessRequiresSuccessfulCheck(t *testing.T) {
	tests := []struct {
		name       string
		check      ReadinessCheck
		wantStatus int
		wantBody   string
	}{
		{
			name:       "missing check",
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   "not ready\n",
		},
		{
			name: "successful check",
			check: func(context.Context) error {
				return nil
			},
			wantStatus: http.StatusOK,
			wantBody:   "ready\n",
		},
		{
			name: "failed check",
			check: func(context.Context) error {
				return errors.New("database password must not be returned")
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   "not ready\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := makeRequest(healthMux(test.check), http.MethodGet, readinessPath, context.Background())

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if response.Body.String() != test.wantBody {
				t.Errorf("body = %q, want %q", response.Body.String(), test.wantBody)
			}
			if strings.Contains(response.Body.String(), "password") {
				t.Error("readiness response exposed the dependency error")
			}
			assertHealthHeaders(t, response)
		})
	}
}

func TestReadinessCheckReceivesRequestContext(t *testing.T) {
	type contextKey struct{}
	ctx := context.WithValue(context.Background(), contextKey{}, "request-context")
	checkCalled := false
	handler := healthMux(func(checkContext context.Context) error {
		checkCalled = true
		if got := checkContext.Value(contextKey{}); got != "request-context" {
			return errors.New("request context was not propagated")
		}
		return nil
	})

	response := makeRequest(handler, http.MethodGet, readinessPath, ctx)

	if !checkCalled {
		t.Fatal("readiness check was not called")
	}
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestHealthEndpointsRejectUnsupportedMethods(t *testing.T) {
	response := makeRequest(healthMux(nil), http.MethodPost, livenessPath, context.Background())

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if got := response.Header().Get("Allow"); got != "GET, HEAD" {
		t.Errorf("Allow = %q, want %q", got, "GET, HEAD")
	}
}

func healthMux(check ReadinessCheck) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterHealthEndpoints(mux, check)
	return mux
}

func makeRequest(handler http.Handler, method, path string, ctx context.Context) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil).WithContext(ctx)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func assertHealthHeaders(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	if got := response.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", got, "text/plain; charset=utf-8")
	}
	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want %q", got, "no-store")
	}
}
