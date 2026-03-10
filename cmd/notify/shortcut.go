package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const notifyConfigPath = "notify.yaml"

// Shortcut represents a saved bus stop monitoring preset.
type Shortcut struct {
	Name         string `yaml:"name"`
	RouteID      string `yaml:"route_id"`
	RouteName    string `yaml:"route_name"`
	StartStop    string `yaml:"start_stop"`
	EndStop      string `yaml:"end_stop"`
	Direction    int    `yaml:"direction"`
	StopID       string `yaml:"stop_id"`
	StopName     string `yaml:"stop_name"`
	StopSequence int    `yaml:"stop_sequence"`
	Threshold    int    `yaml:"threshold"`
}

// NotifyConfig is the top-level structure of notify.yaml.
type NotifyConfig struct {
	Shortcuts []Shortcut `yaml:"shortcuts"`
}

// loadNotifyConfig reads notify.yaml. Returns empty config if file doesn't exist.
func loadNotifyConfig() (*NotifyConfig, error) {
	data, err := os.ReadFile(notifyConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &NotifyConfig{}, nil
		}
		return nil, fmt.Errorf("讀取 %s 失敗: %w", notifyConfigPath, err)
	}

	var cfg NotifyConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 %s 失敗: %w", notifyConfigPath, err)
	}
	return &cfg, nil
}

// saveNotifyConfig writes the config to notify.yaml.
func saveNotifyConfig(cfg *NotifyConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化設定失敗: %w", err)
	}
	if err := os.WriteFile(notifyConfigPath, data, 0644); err != nil {
		return fmt.Errorf("寫入 %s 失敗: %w", notifyConfigPath, err)
	}
	return nil
}

// findShortcut returns the shortcut with the given name, or nil if not found.
func findShortcut(cfg *NotifyConfig, name string) *Shortcut {
	for i := range cfg.Shortcuts {
		if cfg.Shortcuts[i].Name == name {
			return &cfg.Shortcuts[i]
		}
	}
	return nil
}

// formatDirection returns "去程" or "回程".
func formatDirection(direction int) string {
	if direction == 1 {
		return "回程"
	}
	return "去程"
}

// listShortcuts prints all saved shortcuts.
func listShortcuts() {
	cfg, err := loadNotifyConfig()
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}
	if len(cfg.Shortcuts) == 0 {
		fmt.Println("尚無快捷設定")
		return
	}
	fmt.Println("快捷列表：")
	for _, s := range cfg.Shortcuts {
		fmt.Printf("  %-8s %s %s（%s）%d 分鐘\n",
			s.Name, s.RouteName, s.StopName, formatDirection(s.Direction), s.Threshold)
	}
}

// deleteShortcut removes a shortcut by name.
func deleteShortcut(name string) {
	cfg, err := loadNotifyConfig()
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}
	idx := -1
	for i := range cfg.Shortcuts {
		if cfg.Shortcuts[i].Name == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		log.Fatalf("快捷「%s」不存在", name)
	}
	cfg.Shortcuts = append(cfg.Shortcuts[:idx], cfg.Shortcuts[idx+1:]...)
	if err := saveNotifyConfig(cfg); err != nil {
		log.Fatalf("儲存失敗: %v", err)
	}
	fmt.Printf("✓ 已刪除快捷「%s」\n", name)
}

// loadShortcut loads a shortcut by name. Exits on error.
func loadShortcut(name string) *Shortcut {
	cfg, err := loadNotifyConfig()
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}
	s := findShortcut(cfg, name)
	if s == nil {
		names := make([]string, len(cfg.Shortcuts))
		for i, sc := range cfg.Shortcuts {
			names[i] = sc.Name
		}
		if len(names) > 0 {
			log.Fatalf("快捷「%s」不存在\n可用的快捷：%s", name, strings.Join(names, "、"))
		}
		log.Fatalf("快捷「%s」不存在", name)
	}
	fmt.Printf("載入快捷「%s」: %s %s（%s），%d 分鐘前通知\n",
		s.Name, s.RouteName, s.StopName, formatDirection(s.Direction), s.Threshold)
	return s
}
