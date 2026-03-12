package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
	"gopkg.in/yaml.v3"
)

const notifyConfigPath = "notify.yaml"

// Shortcut represents a saved bus stop monitoring preset.
type Shortcut struct {
	Name         string `yaml:"name"`
	RouteName    string `yaml:"route_name"`
	StartStop    string `yaml:"start_stop"`
	EndStop      string `yaml:"end_stop"`
	Direction    int    `yaml:"direction"`
	StopName     string `yaml:"stop_name"`
	StopSequence int    `yaml:"stop_sequence"`
	Threshold    int    `yaml:"threshold"`
	TDXRouteID   string `yaml:"tdx_route_id,omitempty"`
	TDXStopID    string `yaml:"tdx_stop_id,omitempty"`
	EBusRouteID  string `yaml:"ebus_route_id,omitempty"`
	EBusStopID   string `yaml:"ebus_stop_id,omitempty"`
}

// RouteID returns the route ID for the given source ("tdx" or "ebus").
func (s *Shortcut) RouteID(source string) string {
	if source == "tdx" {
		return s.TDXRouteID
	}
	return s.EBusRouteID
}

// StopID returns the stop ID for the given source ("tdx" or "ebus").
func (s *Shortcut) StopID(source string) string {
	if source == "tdx" {
		return s.TDXStopID
	}
	return s.EBusStopID
}

// SetRouteID sets the route ID for the given source.
func (s *Shortcut) SetRouteID(source, id string) {
	if source == "tdx" {
		s.TDXRouteID = id
	} else {
		s.EBusRouteID = id
	}
}

// SetStopID sets the stop ID for the given source.
func (s *Shortcut) SetStopID(source, id string) {
	if source == "tdx" {
		s.TDXStopID = id
	} else {
		s.EBusStopID = id
	}
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

// updateShortcut persists an updated shortcut back to notify.yaml.
func updateShortcut(s *Shortcut) error {
	cfg, err := loadNotifyConfig()
	if err != nil {
		return err
	}
	existing := findShortcut(cfg, s.Name)
	if existing != nil {
		*existing = *s
	}
	return saveNotifyConfig(cfg)
}

// formatDirection returns "去程" or "回程".
func formatDirection(direction int) string {
	if direction == 1 {
		return "回程"
	}
	return "去程"
}

// listShortcuts prints all saved shortcuts with numbered selection.
// Returns the selected shortcut, or nil if cancelled or no shortcuts exist.
func listShortcuts() *Shortcut {
	cfg, err := loadNotifyConfig()
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}
	if len(cfg.Shortcuts) == 0 {
		fmt.Println("尚無快捷設定")
		return nil
	}
	fmt.Println("快捷列表：")
	for i, s := range cfg.Shortcuts {
		fmt.Printf("  [%d] %-8s %s %s（%s）%d 分鐘\n",
			i+1, s.Name, s.RouteName, s.StopName, formatDirection(s.Direction), s.Threshold)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n選擇（Enter 取消）: ")
	if !scanner.Scan() {
		return nil
	}
	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		return nil
	}
	choice, err := strconv.Atoi(text)
	if err != nil || choice < 1 || choice > len(cfg.Shortcuts) {
		fmt.Println("無效選擇")
		return nil
	}
	s := cfg.Shortcuts[choice-1]
	fmt.Printf("載入快捷「%s」: %s %s（%s），%d 分鐘前通知\n",
		s.Name, s.RouteName, s.StopName, formatDirection(s.Direction), s.Threshold)
	return &s
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

// buildShortcut creates a Shortcut with the route/stop IDs stored in the correct
// source-specific fields based on the route's Source.
func buildShortcut(route model.Route, stop model.Stop, direction, threshold int) Shortcut {
	s := Shortcut{
		RouteName:    route.Name,
		StartStop:    route.StartStop,
		EndStop:      route.EndStop,
		Direction:    direction,
		StopName:     stop.Name,
		StopSequence: stop.Sequence,
		Threshold:    threshold,
	}
	s.SetRouteID(route.Source, route.RouteID)
	s.SetStopID(route.Source, stop.StopID)
	return s
}

// promptSaveShortcut asks the user for a shortcut name and saves it.
// If the name already exists, it shows the existing config and asks for overwrite confirmation.
func promptSaveShortcut(scanner *bufio.Scanner, s Shortcut) {
	fmt.Print("\n儲存為快捷？輸入名稱（Enter 跳過）: ")
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		return
	}

	cfg, err := loadNotifyConfig()
	if err != nil {
		log.Printf("載入設定失敗: %v", err)
		return
	}

	existing := findShortcut(cfg, name)
	if existing != nil {
		fmt.Printf("快捷「%s」已存在: %s %s（%s），%d 分鐘\n",
			existing.Name, existing.RouteName, existing.StopName,
			formatDirection(existing.Direction), existing.Threshold)
		fmt.Print("是否覆蓋？(y/N): ")
		if !scanner.Scan() {
			return
		}
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if answer != "y" && answer != "yes" {
			fmt.Println("已取消儲存")
			return
		}
		for i := range cfg.Shortcuts {
			if cfg.Shortcuts[i].Name == name {
				cfg.Shortcuts = append(cfg.Shortcuts[:i], cfg.Shortcuts[i+1:]...)
				break
			}
		}
	}

	s.Name = name
	cfg.Shortcuts = append(cfg.Shortcuts, s)

	if err := saveNotifyConfig(cfg); err != nil {
		log.Printf("儲存失敗: %v", err)
		return
	}
	fmt.Printf("✓ 已儲存快捷「%s」\n", name)
}

// resolveShortcutID uses routeName + stopName to look up the IDs for the given source.
// It searches routes, finds the matching route, gets stops, and matches by stop name.
// Returns the resolved routeID and stopID, or empty strings on failure.
func resolveShortcutID(ctx context.Context, ds model.BusDataSource, city string, s *Shortcut, direction int, source string) (routeID, stopID string, err error) {
	routes, err := ds.SearchRoutes(ctx, city, s.RouteName)
	if err != nil {
		return "", "", fmt.Errorf("反查路線失敗: %w", err)
	}

	var matchedRoute *model.Route
	for i, r := range routes {
		if r.Name == s.RouteName {
			matchedRoute = &routes[i]
			break
		}
	}
	if matchedRoute == nil {
		return "", "", fmt.Errorf("找不到路線「%s」", s.RouteName)
	}

	stops, err := ds.GetStops(ctx, city, matchedRoute.RouteID, direction)
	if err != nil {
		return "", "", fmt.Errorf("反查站點失敗: %w", err)
	}

	for _, st := range stops {
		if st.Name == s.StopName {
			return matchedRoute.RouteID, st.StopID, nil
		}
	}

	return "", "", fmt.Errorf("在路線「%s」中找不到站點「%s」", s.RouteName, s.StopName)
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
