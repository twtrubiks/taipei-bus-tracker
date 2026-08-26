package ebus

import "github.com/twtrubiks/taipei-bus-tracker/internal/model"

// eBus API raw response item (flat array element)
type EBusStopDynRaw struct {
	SN          int       `json:"sn"`          // sequence number (0-based)
	ETA         int       `json:"eta"`         // estimated time of arrival in minutes
	NextDepTime string    `json:"NextDepTime"` // next departure time
	UB          int       `json:"ub"`          // unknown
	BI          []EBusBus `json:"bi"`          // buses approaching this stop (nullable)
	BO          []EBusBus `json:"bo"`          // buses leaving this stop (nullable)
	BISName     string    `json:"bisname"`
}

type EBusBus struct {
	BN    string `json:"bn"` // bus plate number
	BT    int    `json:"bt"`
	BSetL int    `json:"bSetL"`
	BSetN int    `json:"bSetN"`
}

// plates extracts non-empty plate numbers, preserving upstream order.
func plates(buses []EBusBus) []model.Bus {
	var out []model.Bus
	for _, b := range buses {
		if b.BN != "" {
			out = append(out, model.Bus{PlateNumb: b.BN})
		}
	}
	return out
}

func convertETAs(stops []EBusStopDynRaw) []model.StopETA {
	etas := make([]model.StopETA, 0, len(stops))
	for _, s := range stops {
		// Convert ETA from minutes to seconds for unified model
		etaSeconds := s.ETA * 60

		etas = append(etas, model.StopETA{
			StopID:   "",
			StopName: "",       // eBus API doesn't return stop names
			Sequence: s.SN + 1, // eBus is 0-based, unified model is 1-based
			ETA:      etaSeconds,
			Buses:    plates(s.BI),
			// bo = buses that already left this stop, en route to the next one.
			// They stay attached to the stop they left; the UI places them on the
			// segment between this stop and the following one.
			DepartedBuses: plates(s.BO),
			Source:        "ebus",
		})
	}
	return etas
}
