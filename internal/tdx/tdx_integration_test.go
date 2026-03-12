//go:build integration

package tdx

import (
	"context"
	"os"
	"testing"
	"time"
)

func getTDXProvider(t *testing.T) *Provider {
	t.Helper()
	clientID := os.Getenv("TDX_CLIENT_ID")
	clientSecret := os.Getenv("TDX_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("TDX_CLIENT_ID and TDX_CLIENT_SECRET not set")
	}
	return NewProvider(clientID, clientSecret)
}

func TestIntegration_SearchRoutes(t *testing.T) {
	p := getTDXProvider(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routes, err := p.SearchRoutes(ctx, "Taipei", "299")
	if err != nil {
		t.Fatalf("SearchRoutes failed: %v", err)
	}
	if len(routes) == 0 {
		t.Fatal("expected at least 1 route, got 0")
	}
	for _, r := range routes {
		if r.RouteID == "" {
			t.Error("RouteID is empty")
		}
		if r.Name == "" {
			t.Error("Name is empty")
		}
		if r.Source != "tdx" {
			t.Errorf("Source = %q, want tdx", r.Source)
		}
	}
	t.Logf("found %d routes, first: %s (%s)", len(routes), routes[0].Name, routes[0].RouteID)
}

func TestIntegration_GetStops(t *testing.T) {
	p := getTDXProvider(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// First search for a route to get a valid TDX route ID
	routes, err := p.SearchRoutes(ctx, "Taipei", "299")
	if err != nil {
		t.Fatalf("SearchRoutes failed: %v", err)
	}
	if len(routes) == 0 {
		t.Skip("no routes found")
	}

	time.Sleep(1 * time.Second) // rate limit

	stops, err := p.GetStops(ctx, "Taipei", routes[0].RouteID, 0)
	if err != nil {
		t.Fatalf("GetStops failed: %v", err)
	}
	if len(stops) == 0 {
		t.Fatal("expected at least 1 stop, got 0")
	}
	for _, s := range stops {
		if s.StopID == "" {
			t.Error("StopID is empty")
		}
		if s.Name == "" {
			t.Error("Name is empty")
		}
		if s.Sequence <= 0 {
			t.Errorf("Sequence = %d, expected > 0", s.Sequence)
		}
		if s.Source != "tdx" {
			t.Errorf("Source = %q, want tdx", s.Source)
		}
	}
	t.Logf("found %d stops, first: %s (seq %d)", len(stops), stops[0].Name, stops[0].Sequence)
}

func TestIntegration_GetETA(t *testing.T) {
	p := getTDXProvider(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// First search for a route to get a valid TDX route ID
	routes, err := p.SearchRoutes(ctx, "Taipei", "299")
	if err != nil {
		t.Fatalf("SearchRoutes failed: %v", err)
	}
	if len(routes) == 0 {
		t.Skip("no routes found")
	}

	time.Sleep(1 * time.Second) // rate limit

	etas, err := p.GetETA(ctx, "Taipei", routes[0].RouteID, 0)
	if err != nil {
		t.Fatalf("GetETA failed: %v", err)
	}
	if len(etas) == 0 {
		t.Fatal("expected at least 1 ETA, got 0")
	}
	for _, e := range etas {
		if e.Source != "tdx" {
			t.Errorf("Source = %q, want tdx", e.Source)
		}
	}
	t.Logf("found %d ETAs, first: seq %d, eta %d sec", len(etas), etas[0].Sequence, etas[0].ETA)
}
