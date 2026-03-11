package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type TDXConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Config struct {
	Port       int       `yaml:"port"`
	StaticPath string    `yaml:"static_path"`
	Provider   string    `yaml:"provider"`
	TDX        TDXConfig `yaml:"tdx"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Port:       8080,
		StaticPath: "./static",
	}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Environment variables override config file
	if v := os.Getenv("BUS_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Port = p
		}
	}
	if v := os.Getenv("BUS_STATIC_PATH"); v != "" {
		cfg.StaticPath = v
	}
	if v := os.Getenv("TDX_CLIENT_ID"); v != "" {
		cfg.TDX.ClientID = v
	}
	if v := os.Getenv("TDX_CLIENT_SECRET"); v != "" {
		cfg.TDX.ClientSecret = v
	}
	if v := os.Getenv("BUS_PROVIDER"); v != "" {
		cfg.Provider = v
	}

	return cfg, nil
}
