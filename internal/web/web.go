// Package web embeds the built frontend (Vite output) so the backend can serve
// the single-page app from the same binary — no separate frontend image.
//
// internal/web/dist holds a placeholder index.html so `go build` always works;
// the real build is copied into it by the Docker image build before compiling.
package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var dist embed.FS

// FS returns the embedded frontend assets (the contents of dist/).
func FS() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}

// HTTPFS returns the embedded frontend assets as an http.FileSystem.
func HTTPFS() http.FileSystem {
	return http.FS(FS())
}
