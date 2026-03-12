package model

import "context"

// ETA special values (negative = non-arrival status)
const (
	ETANotDeparted   = -1 // 未發車
	ETALastBusLeft   = -2 // 末班車已駛離
	ETANoStop        = -3 // 交管不停靠
	ETANotOperating  = -4 // 未營運
)

type Route struct {
	RouteID   string `json:"routeId"`
	Name      string `json:"routeName"`
	StartStop string `json:"startStop"`
	EndStop   string `json:"endStop"`
	Source    string `json:"source"`
}

type Stop struct {
	StopID   string `json:"stopId"`
	Name     string `json:"stopName"`
	Sequence int    `json:"sequence"`
	Source   string `json:"source"`
}

type Bus struct {
	PlateNumb string `json:"plateNumb"`
}

type StopETA struct {
	StopID   string `json:"stopId"`
	StopName string `json:"stopName"`
	Sequence int    `json:"sequence"`
	ETA      int    `json:"eta"`
	Status   string `json:"status"`
	Buses    []Bus  `json:"buses"`
	Source   string `json:"source"`
}

// BusDataSource defines the interface for bus data providers.
type BusDataSource interface {
	SearchRoutes(ctx context.Context, city, keyword string) ([]Route, error)
	GetStops(ctx context.Context, city, routeID string, direction int) ([]Stop, error)
	GetETA(ctx context.Context, city, routeID string, direction int) ([]StopETA, error)
}
