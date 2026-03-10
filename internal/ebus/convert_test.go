package ebus

import (
	"testing"
)

func TestConvertETAs(t *testing.T) {
	stops := []EBusStopDynRaw{
		{SN: 0, ETA: 5, BI: []EBusBus{{BN: "210-U3", BT: 1}}, BO: nil},
		{SN: 1, ETA: 0, BI: nil, BO: nil},
		{SN: 2, ETA: 3, BI: []EBusBus{{BN: "ABC-12", BT: 1}, {BN: "DEF-34", BT: 1}}, BO: nil},
	}

	etas := convertETAs(stops)
	if len(etas) != 3 {
		t.Fatalf("expected 3 etas, got %d", len(etas))
	}

	// First stop: 5 min = 300 sec, 1 bus
	if etas[0].ETA != 300 {
		t.Errorf("expected eta 300 (5min*60), got %d", etas[0].ETA)
	}
	if etas[0].Sequence != 1 {
		t.Errorf("expected sequence 1 (0-based+1), got %d", etas[0].Sequence)
	}
	if etas[0].Source != "ebus" {
		t.Errorf("expected source ebus, got %s", etas[0].Source)
	}
	if len(etas[0].Buses) != 1 || etas[0].Buses[0].PlateNumb != "210-U3" {
		t.Errorf("expected bus 210-U3, got %v", etas[0].Buses)
	}

	// Second stop: 0 min = 0 sec, no bus
	if etas[1].ETA != 0 {
		t.Errorf("expected eta 0, got %d", etas[1].ETA)
	}
	if len(etas[1].Buses) != 0 {
		t.Errorf("expected 0 buses, got %d", len(etas[1].Buses))
	}

	// Third stop: 2 buses
	if len(etas[2].Buses) != 2 {
		t.Errorf("expected 2 buses, got %d", len(etas[2].Buses))
	}
}

func TestConvertETAs_NilBI(t *testing.T) {
	stops := []EBusStopDynRaw{
		{SN: 0, ETA: 10, BI: nil},
	}
	etas := convertETAs(stops)
	if len(etas[0].Buses) != 0 {
		t.Errorf("expected 0 buses for nil BI, got %d", len(etas[0].Buses))
	}
}

func TestConvertETAs_Empty(t *testing.T) {
	etas := convertETAs(nil)
	if len(etas) != 0 {
		t.Errorf("expected 0 etas, got %d", len(etas))
	}
}
