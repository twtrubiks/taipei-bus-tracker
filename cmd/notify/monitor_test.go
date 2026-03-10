package main

import (
	"testing"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

func TestCalcInterval(t *testing.T) {
	tests := []struct {
		name      string
		etaSec    int
		threshold int // minutes
		want      time.Duration
	}{
		{"far above 2x threshold", 700, 5, 60 * time.Second},
		{"exactly at 2x threshold boundary", 600, 5, 30 * time.Second},
		{"between threshold and 2x", 400, 5, 30 * time.Second},
		{"at threshold", 300, 5, 15 * time.Second},
		{"below threshold", 120, 5, 15 * time.Second},
		{"arriving (0)", 0, 5, 15 * time.Second},
		{"not departed (-1)", -1, 5, 60 * time.Second},
		{"last bus left (-2)", -2, 5, 60 * time.Second},
		{"threshold 3min far", 500, 3, 60 * time.Second},
		{"threshold 3min near", 200, 3, 30 * time.Second},
		{"threshold 3min within", 100, 3, 15 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcInterval(tt.etaSec, tt.threshold)
			if got != tt.want {
				t.Errorf("calcInterval(%d, %d) = %v, want %v", tt.etaSec, tt.threshold, got, tt.want)
			}
		})
	}
}

func TestNotifyState(t *testing.T) {
	tests := []struct {
		name      string
		etas      []int // sequence of ETA values in seconds
		threshold int   // minutes
		wants     []bool
	}{
		{
			name:      "normal arrival",
			etas:      []int{600, 400, 300, 200},
			threshold: 5,
			wants:     []bool{false, false, true, false},
		},
		{
			name:      "two buses sequentially",
			etas:      []int{600, 300, 200, 0, 500, 300, 200},
			threshold: 5,
			wants:     []bool{false, true, false, false, false, true, false},
		},
		{
			name:      "no duplicate notification",
			etas:      []int{300, 250, 200, 150},
			threshold: 5,
			wants:     []bool{true, false, false, false},
		},
		{
			name:      "fluctuation within threshold no retrigger",
			etas:      []int{280, 290, 260, 250},
			threshold: 5,
			wants:     []bool{true, false, false, false},
		},
		{
			name:      "not departed resets state",
			etas:      []int{600, 300, 200, -1, 500, 300},
			threshold: 5,
			wants:     []bool{false, true, false, false, false, true},
		},
		{
			name:      "last bus left resets state",
			etas:      []int{300, 200, -2, 600, 300},
			threshold: 5,
			wants:     []bool{true, false, false, false, true},
		},
		{
			name:      "arriving resets for next bus",
			etas:      []int{300, 200, 0, 400, 300},
			threshold: 5,
			wants:     []bool{true, false, false, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newNotifyState()
			for i, eta := range tt.etas {
				got := s.check(eta, tt.threshold)
				if got != tt.wants[i] {
					t.Errorf("step %d: check(%d, %d) = %v, want %v",
						i, eta, tt.threshold, got, tt.wants[i])
				}
			}
		})
	}
}

func TestFindStopETA(t *testing.T) {
	etas := []model.StopETA{
		{StopID: "S1", Sequence: 1, ETA: 600},
		{StopID: "S2", Sequence: 2, ETA: 300},
		{StopID: "S3", Sequence: 3, ETA: 0},
	}
	tdxETAs := []model.StopETA{
		{StopID: "S1", Sequence: 0, ETA: 600},
		{StopID: "S2", Sequence: 0, ETA: 300},
	}
	tests := []struct {
		name string
		etas []model.StopETA
		stop model.Stop
		want int
	}{
		{"match by sequence", etas, model.Stop{Sequence: 2}, 300},
		{"match by stopId (TDX, sequence=0)", tdxETAs, model.Stop{StopID: "S2", Sequence: 2}, 300},
		{"stopId preferred over sequence", etas, model.Stop{StopID: "S1", Sequence: 3}, 600},
		{"not found returns ETANotDeparted", etas, model.Stop{Sequence: 99}, model.ETANotDeparted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findStopETA(tt.etas, tt.stop)
			if got != tt.want {
				t.Errorf("findStopETA() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFormatETA(t *testing.T) {
	tests := []struct {
		eta  int
		want string
	}{
		{480, "ETA 8 分"},
		{300, "ETA 5 分"},
		{0, "進站中"},
		{-1, "未發車"},
		{-2, "末班駛離"},
		{-3, "交管不停靠"},
		{-4, "未營運"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatETA(tt.eta)
			if got != tt.want {
				t.Errorf("formatETA(%d) = %q, want %q", tt.eta, got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{45 * time.Second, "45s"},
		{5*time.Minute + 30*time.Second, "5m 30s"},
		{2*time.Hour + 15*time.Minute + 10*time.Second, "2h 15m 10s"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatDuration(tt.d)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
