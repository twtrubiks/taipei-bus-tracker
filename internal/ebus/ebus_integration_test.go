//go:build integration

package ebus

import (
	"context"
	"testing"
	"time"
)

const testRouteID = "0100029900" // 299路

func TestIntegration_SearchRoutes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	p := NewProvider()
	routes, err := p.SearchRoutes(ctx, "", "299")
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
		if r.Source != "ebus" {
			t.Errorf("Source = %q, want ebus", r.Source)
		}
	}
	t.Logf("found %d routes, first: %s (%s)", len(routes), routes[0].Name, routes[0].RouteID)
}

func TestIntegration_GetStops(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	p := NewProvider()
	stops, err := p.GetStops(ctx, "", testRouteID, 0)
	if err != nil {
		t.Fatalf("GetStops failed: %v", err)
	}
	if len(stops) == 0 {
		t.Fatal("expected at least 1 stop, got 0")
	}
	for _, s := range stops {
		if s.Name == "" {
			t.Error("stop Name is empty")
		}
		if s.Sequence <= 0 {
			t.Errorf("stop Sequence = %d, expected > 0", s.Sequence)
		}
		if s.Source != "ebus" {
			t.Errorf("Source = %q, want ebus", s.Source)
		}
	}
	t.Logf("found %d stops, first: %s (seq %d)", len(stops), stops[0].Name, stops[0].Sequence)
}

func TestIntegration_GetETA(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	p := NewProvider()
	etas, err := p.GetETA(ctx, "", testRouteID, 0)
	if err != nil {
		t.Fatalf("GetETA failed: %v", err)
	}
	if len(etas) == 0 {
		t.Fatal("expected at least 1 ETA, got 0")
	}
	for _, e := range etas {
		if e.Sequence <= 0 {
			t.Errorf("ETA Sequence = %d, expected > 0", e.Sequence)
		}
		if e.Source != "ebus" {
			t.Errorf("Source = %q, want ebus", e.Source)
		}
	}
	t.Logf("found %d ETAs, first: seq %d, eta %d sec", len(etas), etas[0].Sequence, etas[0].ETA)
}
