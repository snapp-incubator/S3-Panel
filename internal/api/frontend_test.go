package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func serveFrontend(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	e := echo.New()
	e.Use(frontendMiddleware())
	e.GET("/api/ping", func(c echo.Context) error { return c.String(http.StatusTeapot, "api") })

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestFrontendServesIndex(t *testing.T) {
	rec := serveFrontend(t, "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /: got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<html") {
		t.Fatalf("GET / did not return the index.html document")
	}
}

func TestFrontendSPAFallback(t *testing.T) {
	// A client-side route with no matching file must fall back to index.html.
	rec := serveFrontend(t, "/object-storage/s3-bucket/buckets")
	if rec.Code != http.StatusOK {
		t.Fatalf("SPA route: got %d, want 200 (index.html fallback)", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<html") {
		t.Fatalf("SPA route did not fall back to index.html")
	}
}

func TestFrontendSkipsAPI(t *testing.T) {
	// API routes must not be intercepted by the static handler.
	rec := serveFrontend(t, "/api/ping")
	if rec.Code != http.StatusTeapot {
		t.Fatalf("GET /api/ping: got %d, want 418 (handled by API, not static)", rec.Code)
	}
}
