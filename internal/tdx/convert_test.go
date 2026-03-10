package tdx

import (
	"encoding/json"
	"testing"
)

func TestConvertRoutes(t *testing.T) {
	raw := `[{
		"RouteUID": "TPE10001",
		"RouteName": {"Zh_tw": "1"},
		"DepartureStopNameZh": "萬華",
		"DestinationStopNameZh": "松仁路"
	}]`

	var tdxRoutes []TDXRoute
	if err := json.Unmarshal([]byte(raw), &tdxRoutes); err != nil {
		t.Fatal(err)
	}

	routes := convertRoutes(tdxRoutes)
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	r := routes[0]
	if r.RouteID != "TPE10001" {
		t.Errorf("expected routeId TPE10001, got %s", r.RouteID)
	}
	if r.Name != "1" {
		t.Errorf("expected name 1, got %s", r.Name)
	}
	if r.StartStop != "萬華" {
		t.Errorf("expected startStop 萬華, got %s", r.StartStop)
	}
	if r.EndStop != "松仁路" {
		t.Errorf("expected endStop 松仁路, got %s", r.EndStop)
	}
}

func TestConvertRoutes_Empty(t *testing.T) {
	routes := convertRoutes(nil)
	if len(routes) != 0 {
		t.Errorf("expected 0 routes, got %d", len(routes))
	}
}

func TestConvertStops(t *testing.T) {
	raw := `[{
		"Stops": [{
			"StopUID": "TPE10001001",
			"StopName": {"Zh_tw": "萬華車站"},
			"StopSequence": 1
		}, {
			"StopUID": "TPE10001002",
			"StopName": {"Zh_tw": "西門町"},
			"StopSequence": 2
		}]
	}]`

	var tdxStopOfRoutes []TDXStopOfRoute
	if err := json.Unmarshal([]byte(raw), &tdxStopOfRoutes); err != nil {
		t.Fatal(err)
	}

	stops := convertStops(tdxStopOfRoutes)
	if len(stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(stops))
	}
	if stops[0].Name != "萬華車站" {
		t.Errorf("expected 萬華車站, got %s", stops[0].Name)
	}
	if stops[1].Sequence != 2 {
		t.Errorf("expected sequence 2, got %d", stops[1].Sequence)
	}
}

func TestConvertETAs(t *testing.T) {
	raw := `[{
		"StopUID": "TPE10001001",
		"StopName": {"Zh_tw": "萬華車站"},
		"StopSequence": 1,
		"EstimateTime": 300,
		"PlateNumb": "ABC-1234"
	}, {
		"StopUID": "TPE10001002",
		"StopName": {"Zh_tw": "西門町"},
		"StopSequence": 2,
		"EstimateTime": -1
	}]`

	var tdxETAs []TDXETA
	if err := json.Unmarshal([]byte(raw), &tdxETAs); err != nil {
		t.Fatal(err)
	}

	etas := convertETAs(tdxETAs)
	if len(etas) != 2 {
		t.Fatalf("expected 2 etas, got %d", len(etas))
	}

	if etas[0].ETA != 300 {
		t.Errorf("expected eta 300, got %d", etas[0].ETA)
	}
	if etas[0].Source != "tdx" {
		t.Errorf("expected source tdx, got %s", etas[0].Source)
	}
	if len(etas[0].Buses) != 1 || etas[0].Buses[0].PlateNumb != "ABC-1234" {
		t.Errorf("expected bus ABC-1234, got %v", etas[0].Buses)
	}

	if etas[1].ETA != -1 {
		t.Errorf("expected eta -1, got %d", etas[1].ETA)
	}
	if len(etas[1].Buses) != 0 {
		t.Errorf("expected 0 buses, got %d", len(etas[1].Buses))
	}
}

func TestConvertETAs_NoEstimateTime(t *testing.T) {
	raw := `[{
		"StopUID": "TPE10001001",
		"StopName": {"Zh_tw": "萬華車站"},
		"StopSequence": 1,
		"StopStatus": 1
	}]`

	var tdxETAs []TDXETA
	if err := json.Unmarshal([]byte(raw), &tdxETAs); err != nil {
		t.Fatal(err)
	}

	etas := convertETAs(tdxETAs)
	if etas[0].ETA != -1 {
		t.Errorf("expected eta -1 for StopStatus=1, got %d", etas[0].ETA)
	}
}
