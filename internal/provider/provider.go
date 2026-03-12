package provider

import (
	"fmt"

	"github.com/twtrubiks/taipei-bus-tracker/internal/config"
	"github.com/twtrubiks/taipei-bus-tracker/internal/ebus"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
	"github.com/twtrubiks/taipei-bus-tracker/internal/tdx"
)

const (
	ModeAuto = "auto"
	ModeTDX  = "tdx"
	ModeEBus = "ebus"
)

// Build creates primary and fallback BusDataSource based on the provider mode.
// Supported modes: ModeAuto (default), ModeTDX, ModeEBus.
// Returns the normalized mode name along with the providers.
func Build(cfg *config.Config, mode string) (normalizedMode string, primary, fallback model.BusDataSource, err error) {
	if mode == "" {
		mode = ModeAuto
	}

	hasTDX := cfg.TDX.ClientID != "" && cfg.TDX.ClientSecret != ""

	switch mode {
	case ModeAuto:
		if hasTDX {
			primary = tdx.NewProvider(cfg.TDX.ClientID, cfg.TDX.ClientSecret)
			fallback = ebus.NewProvider()
		} else {
			primary = ebus.NewProvider()
		}
	case ModeTDX:
		if !hasTDX {
			return "", nil, nil, fmt.Errorf("provider 模式為 tdx，但缺少 TDX 憑證（設定 tdx.client_id 和 tdx.client_secret）")
		}
		primary = tdx.NewProvider(cfg.TDX.ClientID, cfg.TDX.ClientSecret)
	case ModeEBus:
		primary = ebus.NewProvider()
	default:
		return "", nil, nil, fmt.Errorf("無效的 provider 模式: %q（支援 auto, tdx, ebus）", mode)
	}

	return mode, primary, fallback, nil
}

// PrimarySource returns the source label ("tdx" or "ebus") of the primary provider
// for the given mode and config. This determines which ID fields to use in shortcuts.
func PrimarySource(cfg *config.Config, mode string) string {
	if mode == "" {
		mode = ModeAuto
	}
	switch mode {
	case ModeTDX:
		return ModeTDX
	case ModeEBus:
		return ModeEBus
	default: // ModeAuto
		if cfg.TDX.ClientID != "" && cfg.TDX.ClientSecret != "" {
			return ModeTDX
		}
		return ModeEBus
	}
}
