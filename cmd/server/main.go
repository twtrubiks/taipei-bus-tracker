package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/config"
	"github.com/twtrubiks/taipei-bus-tracker/internal/ebus"
	"github.com/twtrubiks/taipei-bus-tracker/internal/handler"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
	"github.com/twtrubiks/taipei-bus-tracker/internal/tdx"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize providers
	var primary model.BusDataSource
	if cfg.TDX.ClientID != "" && cfg.TDX.ClientSecret != "" {
		primary = tdx.NewProvider(cfg.TDX.ClientID, cfg.TDX.ClientSecret)
		log.Println("TDX provider initialized")
	} else {
		log.Println("WARNING: TDX credentials not set, TDX provider disabled")
	}

	var fallback model.BusDataSource
	ebusProvider := ebus.NewProvider()
	fallback = ebusProvider

	// If no TDX, use eBus as primary
	if primary == nil {
		primary = ebusProvider
		fallback = nil
		log.Println("Using eBus as primary provider (no TDX credentials)")
	}

	c := cache.New(10 * time.Second)
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
