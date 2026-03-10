package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

var errUpstream = fmt.Errorf("upstream error")

type mockProvider struct {
	routes []model.Route
	stops  []model.Stop
	etas   []model.StopETA
	err    error
}

func (m *mockProvider) SearchRoutes(_ context.Context, _, _ string) ([]model.Route, error) {
	return m.routes, m.err
}
func (m *mockProvider) GetStops(_ context.Context, _, _ string, _ int) ([]model.Stop, error) {
	return m.stops, m.err
}
func (m *mockProvider) GetETA(_ context.Context, _, _ string, _ int) ([]model.StopETA, error) {
	return m.etas, m.err
}

func TestFallback_PrimarySuccess(t *testing.T) {
	primary := &mockProvider{etas: []model.StopETA{{StopName: "A", Source: "tdx"}}}
	fallback := &mockProvider{etas: []model.StopETA{{StopName: "A", Source: "ebus"}}}
	c := cache.New(10 * time.Second)

	svc := NewFallbackService(primary, fallback, c)
	etas, err := svc.GetETA(context.Background(), "Taipei", "R1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if etas[0].Source != "tdx" {
		t.Errorf("expected source tdx, got %s", etas[0].Source)
	}
}

func TestFallback_PrimaryFails_FallbackSuccess(t *testing.T) {
	primary := &mockProvider{err: fmt.Errorf("timeout")}
	fallback := &mockProvider{etas: []model.StopETA{{StopName: "A", Source: "ebus"}}}
	c := cache.New(10 * time.Second)

	svc := NewFallbackService(primary, fallback, c)
	etas, err := svc.GetETA(context.Background(), "Taipei", "R1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if etas[0].Source != "ebus" {
		t.Errorf("expected source ebus, got %s", etas[0].Source)
	}
}

func TestFallback_BothFail_CacheHit(t *testing.T) {
	primary := &mockProvider{err: fmt.Errorf("timeout")}
	fallback := &mockProvider{err: fmt.Errorf("csrf error")}
	c := cache.New(10 * time.Second)
	c.Set("eta:Taipei:R1:0", []model.StopETA{{StopName: "A", Source: "cached"}})

	svc := NewFallbackService(primary, fallback, c)
	etas, err := svc.GetETA(context.Background(), "Taipei", "R1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if etas[0].Source != "cached" {
		t.Errorf("expected source cached, got %s", etas[0].Source)
	}
}

func TestFallback_BothFail_NoCacheHit(t *testing.T) {
	primary := &mockProvider{err: fmt.Errorf("timeout")}
	fallback := &mockProvider{err: fmt.Errorf("csrf error")}
	c := cache.New(10 * time.Second)

	svc := NewFallbackService(primary, fallback, c)
	_, err := svc.GetETA(context.Background(), "Taipei", "R1", 0)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFallback_PrimaryFails_NoFallback_CacheHit(t *testing.T) {
	primary := &mockProvider{err: fmt.Errorf("timeout")}
	c := cache.New(10 * time.Second)
	c.Set("eta:Taipei:R1:0", []model.StopETA{{StopName: "A", Source: "cached"}})

	svc := NewFallbackService(primary, nil, c)
	etas, err := svc.GetETA(context.Background(), "Taipei", "R1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if etas[0].Source != "cached" {
		t.Errorf("expected source cached, got %s", etas[0].Source)
	}
}

func TestFallback_PrimaryFails_NoFallback_NoCacheHit(t *testing.T) {
	primary := &mockProvider{err: fmt.Errorf("timeout")}
	c := cache.New(10 * time.Second)

	svc := NewFallbackService(primary, nil, c)
	_, err := svc.GetETA(context.Background(), "Taipei", "R1", 0)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
