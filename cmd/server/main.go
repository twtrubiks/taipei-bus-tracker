package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received %v, shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on %s (static: %s)", addr, cfg.StaticPath)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
	log.Println("Server stopped")
}
