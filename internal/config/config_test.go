package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("nonexistent.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.StaticPath != "./static" {
		t.Errorf("expected static_path ./static, got %s", cfg.StaticPath)
	}
}

func TestLoad_YAMLFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte("port: 9090\nstatic_path: /var/www\ntdx:\n  client_id: myid\n  client_secret: mysecret\n")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.TDX.ClientID != "myid" {
		t.Errorf("expected client_id myid, got %s", cfg.TDX.ClientID)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte("port: 8080\ntdx:\n  client_id: file_id\n")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("BUS_PORT", "3000")
	t.Setenv("TDX_CLIENT_ID", "env_id")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Port)
	}
	if cfg.TDX.ClientID != "env_id" {
		t.Errorf("expected client_id env_id, got %s", cfg.TDX.ClientID)
	}
}
