package model

import "fmt"

// ETAStatus converts an ETA value (in seconds) to a human-readable status string.
func ETAStatus(eta int) string {
	switch {
	case eta == ETANotDeparted:
		return "未發車"
	case eta == ETALastBusLeft:
		return "末班車已駛離"
	case eta == ETANoStop:
		return "交管不停靠"
	case eta == ETANotOperating:
		return "未營運"
	case eta >= 0 && eta <= 180:
		return "進站中"
	case eta > 0:
		minutes := (eta + 59) / 60 // round up
		return fmt.Sprintf("約%d分", minutes)
	default:
		return "未知"
	}
}
