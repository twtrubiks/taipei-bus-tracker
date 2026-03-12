package tdx

import "github.com/twtrubiks/taipei-bus-tracker/internal/model"

// TDX API response structures

type NameType struct {
	ZhTw string `json:"Zh_tw"`
}

type TDXRoute struct {
	RouteUID             string   `json:"RouteUID"`
	RouteName            NameType `json:"RouteName"`
	DepartureStopNameZh  string   `json:"DepartureStopNameZh"`
	DestinationStopNameZh string  `json:"DestinationStopNameZh"`
}

type TDXStop struct {
	StopUID      string   `json:"StopUID"`
	StopName     NameType `json:"StopName"`
	StopSequence int      `json:"StopSequence"`
}

type TDXStopOfRoute struct {
	Stops []TDXStop `json:"Stops"`
}

type TDXETA struct {
	StopUID      string   `json:"StopUID"`
	StopName     NameType `json:"StopName"`
	StopSequence int      `json:"StopSequence"`
	EstimateTime *int     `json:"EstimateTime"`
	StopStatus   int      `json:"StopStatus"`
	PlateNumb    string   `json:"PlateNumb"`
}

func convertRoutes(tdxRoutes []TDXRoute) []model.Route {
	routes := make([]model.Route, 0, len(tdxRoutes))
	for _, r := range tdxRoutes {
		routes = append(routes, model.Route{
			RouteID:   r.RouteUID,
			Name:      r.RouteName.ZhTw,
			StartStop: r.DepartureStopNameZh,
			EndStop:   r.DestinationStopNameZh,
			Source:    "tdx",
		})
	}
	return routes
}

func convertStops(stopOfRoutes []TDXStopOfRoute) []model.Stop {
	// Use only the first StopOfRoute record to avoid duplicates from multiple sub-routes
	if len(stopOfRoutes) == 0 {
		return nil
	}
	first := stopOfRoutes[0]
	stops := make([]model.Stop, 0, len(first.Stops))
	for _, s := range first.Stops {
		stops = append(stops, model.Stop{
			StopID:   s.StopUID,
			Name:     s.StopName.ZhTw,
			Sequence: s.StopSequence,
			Source:   "tdx",
		})
	}
	return stops
}

func convertETAs(tdxETAs []TDXETA) []model.StopETA {
	etas := make([]model.StopETA, 0, len(tdxETAs))
	for _, e := range tdxETAs {
		eta := etaFromTDX(e)
		var buses []model.Bus
		if e.PlateNumb != "" {
			buses = append(buses, model.Bus{PlateNumb: e.PlateNumb})
		}
		etas = append(etas, model.StopETA{
			StopID:   e.StopUID,
			StopName: e.StopName.ZhTw,
			Sequence: e.StopSequence,
			ETA:      eta,
			Buses:    buses,
			Source:   "tdx",
		})
	}
	return etas
}

// etaFromTDX converts TDX EstimateTime and StopStatus to unified ETA value.
// StopStatus: 0=normal, 1=not departed, 2=route ended, 3=not operating, 4=detour
func etaFromTDX(e TDXETA) int {
	if e.EstimateTime != nil {
		return *e.EstimateTime
	}
	switch e.StopStatus {
	case 1:
		return model.ETANotDeparted
	case 2:
		return model.ETALastBusLeft
	case 3:
		return model.ETANotOperating
	case 4:
		return model.ETANoStop
	default:
		return model.ETANotDeparted
	}
}
