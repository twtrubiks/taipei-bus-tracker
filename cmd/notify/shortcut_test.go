package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

// mockDS is a mock BusDataSource for testing lazy resolve.
type mockDS struct {
	routes []model.Route
	stops  []model.Stop
	err    error
}

func (m *mockDS) SearchRoutes(_ context.Context, _, _ string) ([]model.Route, error) {
	return m.routes, m.err
}
func (m *mockDS) GetStops(_ context.Context, _, _ string, _ int) ([]model.Stop, error) {
	return m.stops, m.err
}
func (m *mockDS) GetETA(_ context.Context, _, _ string, _ int) ([]model.StopETA, error) {
	return nil, nil
}

func TestResolveShortcutID_Success(t *testing.T) {
	ds := &mockDS{
		routes: []model.Route{
			{RouteID: "EB123", Name: "299", Source: "ebus"},
		},
		stops: []model.Stop{
			{StopID: "S1", Name: "輔大", Sequence: 1, Source: "ebus"},
			{StopID: "S2", Name: "台北車站", Sequence: 5, Source: "ebus"},
		},
	}

	s := &Shortcut{
		RouteName: "299",
		StopName:  "台北車站",
	}

	routeID, stopID, err := resolveShortcutID(context.Background(), ds, "Taipei", s, 0, "ebus")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if routeID != "EB123" {
		t.Errorf("routeID = %q, want EB123", routeID)
	}
	if stopID != "S2" {
		t.Errorf("stopID = %q, want S2", stopID)
	}
}

func TestResolveShortcutID_RouteNotFound(t *testing.T) {
	ds := &mockDS{
		routes: []model.Route{
			{RouteID: "EB123", Name: "300", Source: "ebus"},
		},
	}

	s := &Shortcut{
		RouteName: "299",
		StopName:  "台北車站",
	}

	_, _, err := resolveShortcutID(context.Background(), ds, "Taipei", s, 0, "ebus")
	if err == nil {
		t.Error("expected error when route not found")
	}
}

func TestResolveShortcutID_StopNotFound(t *testing.T) {
	ds := &mockDS{
		routes: []model.Route{
			{RouteID: "EB123", Name: "299", Source: "ebus"},
		},
		stops: []model.Stop{
			{StopID: "S1", Name: "輔大", Sequence: 1},
		},
	}

	s := &Shortcut{
		RouteName: "299",
		StopName:  "不存在的站",
	}

	_, _, err := resolveShortcutID(context.Background(), ds, "Taipei", s, 0, "ebus")
	if err == nil {
		t.Error("expected error when stop not found")
	}
}

func TestResolveShortcutID_SearchFails(t *testing.T) {
	ds := &mockDS{err: fmt.Errorf("network error")}

	s := &Shortcut{
		RouteName: "299",
		StopName:  "台北車站",
	}

	_, _, err := resolveShortcutID(context.Background(), ds, "Taipei", s, 0, "ebus")
	if err == nil {
		t.Error("expected error on search failure")
	}
}

func TestShortcutHelpers(t *testing.T) {
	s := Shortcut{}

	s.SetRouteID("tdx", "TPE123")
	s.SetStopID("tdx", "TPE456")
	s.SetRouteID("ebus", "EB789")
	s.SetStopID("ebus", "EB012")

	if s.RouteID("tdx") != "TPE123" {
		t.Errorf("RouteID(tdx) = %q, want TPE123", s.RouteID("tdx"))
	}
	if s.StopID("tdx") != "TPE456" {
		t.Errorf("StopID(tdx) = %q, want TPE456", s.StopID("tdx"))
	}
	if s.RouteID("ebus") != "EB789" {
		t.Errorf("RouteID(ebus) = %q, want EB789", s.RouteID("ebus"))
	}
	if s.StopID("ebus") != "EB012" {
		t.Errorf("StopID(ebus) = %q, want EB012", s.StopID("ebus"))
	}
}
