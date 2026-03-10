package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

// SPAHandler serves static files and falls back to index.html for SPA routing.
func SPAHandler(staticPath string) http.Handler {
	fs := http.Dir(staticPath)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if file exists
		path := filepath.Join(staticPath, r.URL.Path)
		if _, err := os.Stat(path); err == nil {
			http.FileServer(fs).ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
	})
}
