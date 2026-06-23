package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestProvideRegionEndpoints ensures the [server.region_endpoints] sub-table
// (with a dashed key) round-trips through koanf into the map.
func TestProvideRegionEndpoints(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	content := `
[server]
region = "teh-1"
serve_frontend = true

[server.region_endpoints]
"teh-2" = "http://s3-panel.apps.private.okd4.teh-2.snappcloud.io"

[object_storage_config]
url = "https://s3.teh-1.snappcloud.io"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := Provide(path)

	if cfg.Server.Region != "teh-1" {
		t.Fatalf("region: got %q, want teh-1", cfg.Server.Region)
	}
	want := "http://s3-panel.apps.private.okd4.teh-2.snappcloud.io"
	if got := cfg.Server.RegionEndpoints["teh-2"]; got != want {
		t.Fatalf("region_endpoints[teh-2]: got %q, want %q", got, want)
	}
	if !cfg.Server.ServeFrontend {
		t.Fatalf("serve_frontend: got false, want true")
	}
}
