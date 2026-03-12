package provider

import (
	"testing"

	"github.com/twtrubiks/taipei-bus-tracker/internal/config"
)

func TestBuild_AutoWithTDXKey(t *testing.T) {
	cfg := &config.Config{
		TDX: config.TDXConfig{ClientID: "id", ClientSecret: "secret"},
	}
	mode, primary, fallback, err := Build(cfg, "auto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "auto" {
		t.Errorf("mode = %q, want auto", mode)
	}
	if primary == nil {
		t.Error("expected primary provider, got nil")
	}
	if fallback == nil {
		t.Error("expected fallback provider, got nil")
	}
}

func TestBuild_AutoWithoutTDXKey(t *testing.T) {
	cfg := &config.Config{}
	mode, primary, fallback, err := Build(cfg, "auto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "auto" {
		t.Errorf("mode = %q, want auto", mode)
	}
	if primary == nil {
		t.Error("expected primary provider (ebus), got nil")
	}
	if fallback != nil {
		t.Error("expected no fallback, got non-nil")
	}
}

func TestBuild_TDXMode(t *testing.T) {
	cfg := &config.Config{
		TDX: config.TDXConfig{ClientID: "id", ClientSecret: "secret"},
	}
	mode, primary, fallback, err := Build(cfg, "tdx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "tdx" {
		t.Errorf("mode = %q, want tdx", mode)
	}
	if primary == nil {
		t.Error("expected primary provider, got nil")
	}
	if fallback != nil {
		t.Error("expected no fallback for tdx mode, got non-nil")
	}
}

func TestBuild_EBusMode(t *testing.T) {
	cfg := &config.Config{}
	mode, primary, fallback, err := Build(cfg, "ebus")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "ebus" {
		t.Errorf("mode = %q, want ebus", mode)
	}
	if primary == nil {
		t.Error("expected primary provider, got nil")
	}
	if fallback != nil {
		t.Error("expected no fallback for ebus mode, got non-nil")
	}
}

func TestBuild_TDXMissingCredentials(t *testing.T) {
	cfg := &config.Config{}
	_, _, _, err := Build(cfg, "tdx")
	if err == nil {
		t.Error("expected error for tdx mode without credentials")
	}
}

func TestBuild_InvalidMode(t *testing.T) {
	cfg := &config.Config{}
	_, _, _, err := Build(cfg, "invalid")
	if err == nil {
		t.Error("expected error for invalid mode")
	}
}

func TestPrimarySource(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		mode string
		want string
	}{
		{"auto with TDX", &config.Config{TDX: config.TDXConfig{ClientID: "id", ClientSecret: "s"}}, "auto", "tdx"},
		{"auto without TDX", &config.Config{}, "auto", "ebus"},
		{"tdx mode", &config.Config{}, "tdx", "tdx"},
		{"ebus mode", &config.Config{}, "ebus", "ebus"},
		{"empty defaults to auto+ebus", &config.Config{}, "", "ebus"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrimarySource(tt.cfg, tt.mode)
			if got != tt.want {
				t.Errorf("PrimarySource() = %q, want %q", got, tt.want)
			}
		})
	}
}
