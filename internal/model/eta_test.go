package model

import "testing"

func TestETAStatus(t *testing.T) {
	tests := []struct {
		eta      int
		expected string
	}{
		{0, "進站中"},
		{60, "進站中"},
		{ArrivedMaxSec, "進站中"},
		{ArrivedMaxSec + 1, "將到站"},
		{179, "將到站"},
		{ArrivingMaxSec, "約3分"},
		{181, "約4分"}, // rounds up
		{300, "約5分"},
		{600, "約10分"},
		{ETANotDeparted, "未發車"},
		{ETALastBusLeft, "末班車已駛離"},
		{ETANoStop, "交管不停靠"},
		{ETANotOperating, "未營運"},
		{-99, "未知"},
	}

	for _, tt := range tests {
		got := ETAStatus(tt.eta)
		if got != tt.expected {
			t.Errorf("ETAStatus(%d) = %q, want %q", tt.eta, got, tt.expected)
		}
	}
}
