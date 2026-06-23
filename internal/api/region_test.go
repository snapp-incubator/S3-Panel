package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/snapp-incubator/S3-Panel/internal/config"
)

func newRegionServer(t *testing.T, local string, endpoints map[string]string) *Server {
	t.Helper()
	s := &Server{
		Config: config.Config{
			Server: config.ServerConfig{Region: local, RegionEndpoints: endpoints},
		},
		logger: zap.NewNop(),
	}
	if err := s.buildRegionTargets(); err != nil {
		t.Fatalf("buildRegionTargets: %v", err)
	}
	return s
}

func regionEcho(s *Server) *echo.Echo {
	e := echo.New()
	g := e.Group("/api/bucket", s.regionRouter())
	g.GET("/list", func(c echo.Context) error { return c.String(http.StatusOK, "local") })
	e.GET("/api/regions", s.HandleRegions())
	return e
}

func doGet(e *echo.Echo, path, region string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if region != "" {
		req.Header.Set("region", region)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestRegionHandledLocally(t *testing.T) {
	e := regionEcho(newRegionServer(t, "teh-1", nil))
	for _, region := range []string{"", "teh-1"} {
		rec := doGet(e, "/api/bucket/list", region)
		if rec.Code != http.StatusOK || rec.Body.String() != "local" {
			t.Fatalf("region %q: got %d %q, want 200 local", region, rec.Code, rec.Body.String())
		}
	}
}

func TestRegionProxied(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("remote:" + r.URL.Path))
	}))
	defer upstream.Close()

	e := regionEcho(newRegionServer(t, "teh-1", map[string]string{"teh-2": upstream.URL}))

	rec := doGet(e, "/api/bucket/list", "teh-2")
	if rec.Code != http.StatusOK || rec.Body.String() != "remote:/api/bucket/list" {
		t.Fatalf("proxy: got %d %q, want 200 remote:/api/bucket/list", rec.Code, rec.Body.String())
	}
}

func TestRegionUnknownRejected(t *testing.T) {
	e := regionEcho(newRegionServer(t, "teh-1", nil))
	rec := doGet(e, "/api/bucket/list", "mars")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown region: got %d, want 400", rec.Code)
	}
}

func TestRegionsListed(t *testing.T) {
	e := regionEcho(newRegionServer(t, "teh-1", map[string]string{"teh-2": "http://teh-2"}))
	rec := doGet(e, "/api/regions", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("regions: got %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "teh-1") || !strings.Contains(body, "teh-2") {
		t.Fatalf("regions body missing entries: %s", body)
	}
}
