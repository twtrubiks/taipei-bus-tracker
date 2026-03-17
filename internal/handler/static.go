package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SPAHandler serves static files and falls back to index.html for SPA routing.
func SPAHandler(staticPath string) http.Handler {
	staticPath = filepath.Clean(staticPath)
	fs := http.Dir(staticPath)
	fileServer := http.FileServer(fs)
	indexPath := filepath.Join(staticPath, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if file exists (clean path to prevent traversal)
		path := filepath.Join(staticPath, filepath.Clean("/"+r.URL.Path))
		if !strings.HasPrefix(path, staticPath) {
			http.ServeFile(w, r, indexPath)
			return
		}
		if _, err := os.Stat(path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		http.ServeFile(w, r, indexPath)
	})
}
