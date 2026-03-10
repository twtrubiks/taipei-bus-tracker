package model

import "testing"

func TestETAStatus(t *testing.T) {
	tests := []struct {
		eta      int
		expected string
	}{
		{300, "約5分"},
		{60, "進站中"},
		{180, "進站中"},
		{0, "進站中"},
		{-1, "未發車"},
		{-2, "末班車已駛離"},
		{-3, "交管不停靠"},
		{-4, "未營運"},
		{600, "約10分"},
		{61, "進站中"},
		{181, "約4分"}, // 181 sec -> rounds up to 4 min
	}

	for _, tt := range tests {
		got := ETAStatus(tt.eta)
		if got != tt.expected {
			t.Errorf("ETAStatus(%d) = %q, want %q", tt.eta, got, tt.expected)
		}
	}
}
