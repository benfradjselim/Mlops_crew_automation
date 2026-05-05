package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

// Handler returns an http.Handler that serves the embedded dashboard at /ui.
func Handler() http.Handler {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic("ui: embed fs sub failed: " + err.Error())
	}
	return http.StripPrefix("/ui", http.FileServer(http.FS(sub)))
}
