package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/config"
	"github.com/twtrubiks/taipei-bus-tracker/internal/handler"
	"github.com/twtrubiks/taipei-bus-tracker/internal/provider"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize providers based on config provider mode
	mode, primary, fallback, err := provider.Build(cfg, cfg.Provider)
	if err != nil {
		log.Fatalf("provider init failed: %v", err)
	}
	log.Printf("Provider mode: %s (primary: %s)", mode, provider.PrimarySource(cfg, mode))

	c := cache.New(30 * time.Second)
	defer c.Close()

	h := handler.NewHandlers(primary, fallback, c)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/routes/search", h.SearchRoutes)
	mux.HandleFunc("GET /api/routes/{routeId}/stops", h.GetStops)
	mux.HandleFunc("GET /api/routes/{routeId}/eta", h.GetETA)
	mux.Handle("/", handler.SPAHandler(cfg.StaticPath))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s (static: %s)", addr, cfg.StaticPath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
