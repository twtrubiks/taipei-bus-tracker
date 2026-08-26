package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

func setupHandlers(primary *mockProvider, fallback *mockProvider) *Handlers {
	c := cache.New(10 * time.Second)
	var fb model.BusDataSource
	if fallback != nil {
		fb = fallback
	}
	return NewHandlers(primary, fb, c)
}

func TestSearchRoutes_Success(t *testing.T) {
	primary := &mockProvider{
		routes: []model.Route{{RouteID: "R1", Name: "1", StartStop: "A", EndStop: "B"}},
	}
	h := setupHandlers(primary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/search?q=1", nil)
	w := httptest.NewRecorder()
	h.SearchRoutes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var routes []model.Route
	if err := json.NewDecoder(w.Body).Decode(&routes); err != nil {
		t.Fatal(err)
	}
	if len(routes) != 1 || routes[0].Name != "1" {
		t.Errorf("unexpected routes: %v", routes)
	}
}

func TestSearchRoutes_MissingQ(t *testing.T) {
	h := setupHandlers(&mockProvider{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/search", nil)
	w := httptest.NewRecorder()
	h.SearchRoutes(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSearchRoutes_Empty(t *testing.T) {
	primary := &mockProvider{routes: []model.Route{}}
	h := setupHandlers(primary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/search?q=nonexistent", nil)
	w := httptest.NewRecorder()
	h.SearchRoutes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetETA_Success(t *testing.T) {
	primary := &mockProvider{
		etas: []model.StopETA{
			{StopName: "站A", Sequence: 1, ETA: 300, Source: "tdx"},
			{StopName: "站B", Sequence: 2, ETA: -1, Source: "tdx"},
		},
	}
	h := setupHandlers(primary, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/routes/{routeId}/eta", h.GetETA)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/R1/eta?gb=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp ETAResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Source != "tdx" {
		t.Errorf("expected source tdx, got %s", resp.Source)
	}
	if len(resp.Stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(resp.Stops))
	}
	if resp.Stops[0].ETA != 300 {
		t.Errorf("expected eta 300, got %d", resp.Stops[0].ETA)
	}
	if resp.Stops[1].ETA != model.ETANotDeparted {
		t.Errorf("expected eta %d, got %d", model.ETANotDeparted, resp.Stops[1].ETA)
	}
}

// The ETA response must carry departedBuses through to the client, attached to
// the stop the bus left rather than the one it is heading to.
func TestGetETA_DepartedBuses(t *testing.T) {
	primary := &mockProvider{
		etas: []model.StopETA{
			{StopName: "站A", Sequence: 1, ETA: 180, Source: "ebus",
				DepartedBuses: []model.Bus{{PlateNumb: "KKB-1789"}}},
			{StopName: "站B", Sequence: 2, ETA: 60, Source: "ebus",
				Buses: []model.Bus{{PlateNumb: "EAL-5812"}}},
		},
	}
	h := setupHandlers(primary, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/routes/{routeId}/eta", h.GetETA)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/R1/eta?gb=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.Bytes()
	if !bytes.Contains(body, []byte(`"departedBuses"`)) {
		t.Errorf("response must include departedBuses field, got %s", body)
	}

	var resp ETAResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Stops[0].DepartedBuses) != 1 ||
		resp.Stops[0].DepartedBuses[0].PlateNumb != "KKB-1789" {
		t.Errorf("stop 1 departedBuses = %v, want [KKB-1789]", resp.Stops[0].DepartedBuses)
	}
	if len(resp.Stops[0].Buses) != 0 {
		t.Errorf("stop 1 buses = %v, want empty", resp.Stops[0].Buses)
	}
	// The departed bus must not be moved onto the following stop.
	if len(resp.Stops[1].DepartedBuses) != 0 {
		t.Errorf("stop 2 departedBuses = %v, want empty", resp.Stops[1].DepartedBuses)
	}
	if len(resp.Stops[1].Buses) != 1 ||
		resp.Stops[1].Buses[0].PlateNumb != "EAL-5812" {
		t.Errorf("stop 2 buses = %v, want [EAL-5812]", resp.Stops[1].Buses)
	}
}

func TestGetETA_EBus_FillsStopInfo(t *testing.T) {
	primary := &mockProvider{
		etas: []model.StopETA{
			{Sequence: 1, ETA: 300, Source: "ebus"},
			{Sequence: 2, ETA: -1, Source: "ebus"},
		},
		stops: []model.Stop{
			{StopID: "S1", Name: "建國中學", Sequence: 1},
			{StopID: "S2", Name: "台大醫院", Sequence: 2},
		},
	}
	h := setupHandlers(primary, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/routes/{routeId}/eta", h.GetETA)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/R1/eta?gb=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp ETAResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Source != "ebus" {
		t.Errorf("expected source ebus, got %s", resp.Source)
	}
	if len(resp.Stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(resp.Stops))
	}
	if resp.Stops[0].StopID != "S1" {
		t.Errorf("expected stopId S1, got %s", resp.Stops[0].StopID)
	}
	if resp.Stops[0].StopName != "建國中學" {
		t.Errorf("expected stopName 建國中學, got %s", resp.Stops[0].StopName)
	}
	if resp.Stops[0].ETA != 300 {
		t.Errorf("expected eta 300, got %d", resp.Stops[0].ETA)
	}
	if resp.Stops[1].StopID != "S2" {
		t.Errorf("expected stopId S2, got %s", resp.Stops[1].StopID)
	}
	if resp.Stops[1].StopName != "台大醫院" {
		t.Errorf("expected stopName 台大醫院, got %s", resp.Stops[1].StopName)
	}
}

func TestGetETA_UpstreamError(t *testing.T) {
	primary := &mockProvider{err: errUpstream}
	h := setupHandlers(primary, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/routes/{routeId}/eta", h.GetETA)

	req := httptest.NewRequest(http.MethodGet, "/api/routes/R1/eta?gb=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected 502, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["error"] != "upstream unavailable" {
		t.Errorf("expected 'upstream unavailable', got %s", resp["error"])
	}
}
