package model

import "fmt"

// ETA tier boundaries in seconds. Keep in sync with web/src/utils/etaStatus.ts —
// both sides are covered by the same boundary cases.
const (
	// ArrivedMaxSec is the upper bound of the "pulling in" tier.
	ArrivedMaxSec = 90
	// ArrivingMaxSec is the exclusive upper bound of the "about to arrive" tier.
	ArrivingMaxSec = 180
)

// ETAStatus converts an ETA value (in seconds) to a human-readable status string.
func ETAStatus(eta int) string {
	switch eta {
	case ETANotDeparted:
		return "未發車"
	case ETALastBusLeft:
		return "末班車已駛離"
	case ETANoStop:
		return "交管不停靠"
	case ETANotOperating:
		return "未營運"
	}
	switch {
	case eta < 0:
		return "未知"
	case eta <= ArrivedMaxSec:
		return "進站中"
	case eta < ArrivingMaxSec:
		return "將到站"
	default:
		minutes := (eta + 59) / 60 // round up
		return fmt.Sprintf("約%d分", minutes)
	}
}
