package ebus

import (
	"slices"
	"testing"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
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

// plateList flattens buses to plate numbers for comparison.
func plateList(buses []model.Bus) []string {
	out := make([]string, 0, len(buses))
	for _, b := range buses {
		out = append(out, b.PlateNumb)
	}
	return out
}

// bo (buses that left a stop) must survive conversion — dropping it was the
// cause of plate numbers vanishing while a bus was between two stops.
func TestConvertETAs_DepartedBuses(t *testing.T) {
	stops := []EBusStopDynRaw{
		// only bi: bus is standing at the stop
		{SN: 0, ETA: 0, BI: []EBusBus{{BN: "EAL-1293"}}, BO: nil},
		// only bo: bus already left, en route to the next stop
		{SN: 1, ETA: 3, BI: nil, BO: []EBusBus{{BN: "KKB-1789"}}},
		// both: one standing, two already gone
		{SN: 2, ETA: 1, BI: []EBusBus{{BN: "EAL-5812"}}, BO: []EBusBus{{BN: "KKA-0161"}, {BN: "KKB-1785"}}},
		// neither
		{SN: 3, ETA: 7, BI: nil, BO: nil},
	}

	etas := convertETAs(stops)

	cases := []struct {
		name     string
		idx      int
		buses    []string
		departed []string
	}{
		{"bi only", 0, []string{"EAL-1293"}, nil},
		{"bo only", 1, nil, []string{"KKB-1789"}},
		{"bi and bo", 2, []string{"EAL-5812"}, []string{"KKA-0161", "KKB-1785"}},
		{"neither", 3, nil, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// slices.Equal treats nil and empty as equal, which is what we want
			// here: "no buses" may arrive as either.
			if got := plateList(etas[c.idx].Buses); !slices.Equal(got, c.buses) {
				t.Errorf("Buses = %v, want %v", got, c.buses)
			}
			if got := plateList(etas[c.idx].DepartedBuses); !slices.Equal(got, c.departed) {
				t.Errorf("DepartedBuses = %v, want %v", got, c.departed)
			}
		})
	}
}

// A departed bus stays attached to the stop it left, not the stop it is heading to.
func TestConvertETAs_DepartedBusStaysOnSourceStop(t *testing.T) {
	stops := []EBusStopDynRaw{
		{SN: 12, ETA: 3, BI: nil, BO: []EBusBus{{BN: "KKB-1789"}}},
		{SN: 13, ETA: 1, BI: nil, BO: nil},
	}
	etas := convertETAs(stops)

	if len(etas[0].DepartedBuses) != 1 || etas[0].DepartedBuses[0].PlateNumb != "KKB-1789" {
		t.Errorf("stop sequence %d should carry the departed bus, got %v", etas[0].Sequence, etas[0].DepartedBuses)
	}
	if len(etas[1].Buses) != 0 || len(etas[1].DepartedBuses) != 0 {
		t.Errorf("next stop must not receive the departed bus, got buses=%v departed=%v",
			etas[1].Buses, etas[1].DepartedBuses)
	}
}

// Upstream sometimes includes entries without a plate number; they carry no
// information and must not become blank labels in the UI.
func TestConvertETAs_SkipsEmptyPlates(t *testing.T) {
	stops := []EBusStopDynRaw{
		{SN: 0, ETA: 2, BI: []EBusBus{{BN: ""}}, BO: []EBusBus{{BN: ""}, {BN: "171-U7"}}},
	}
	etas := convertETAs(stops)

	if len(etas[0].Buses) != 0 {
		t.Errorf("expected empty plate to be skipped, got %v", etas[0].Buses)
	}
	if got := plateList(etas[0].DepartedBuses); len(got) != 1 || got[0] != "171-U7" {
		t.Errorf("DepartedBuses = %v, want [171-U7]", got)
	}
}
