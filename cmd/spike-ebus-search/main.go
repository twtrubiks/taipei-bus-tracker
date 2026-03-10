package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/ebus"
)

func main() {
	keyword := "299"
	if len(os.Args) > 1 {
		keyword = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	p := ebus.NewProvider()

	// Step 1: SearchRoutes
	fmt.Printf("=== Step 1: SearchRoutes('%s') ===\n", keyword)
	routes, err := p.SearchRoutes(ctx, "", keyword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SearchRoutes 失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("找到 %d 條路線\n", len(routes))
	for i, r := range routes {
		fmt.Printf("  [%d] %s (ID: %s) %s ↔ %s\n", i+1, r.Name, r.RouteID, r.StartStop, r.EndStop)
	}

	if len(routes) == 0 {
		fmt.Println("無路線，結束")
		return
	}

	// Step 2: GetStops (使用第一條路線)
	route := routes[0]
	fmt.Printf("\n=== Step 2: GetStops('%s', direction=0) ===\n", route.RouteID)
	stops, err := p.GetStops(ctx, "", route.RouteID, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetStops 失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("找到 %d 個站點（去程）\n", len(stops))
	limit := 5
	if len(stops) < limit {
		limit = len(stops)
	}
	for i := 0; i < limit; i++ {
		s := stops[i]
		fmt.Printf("  站序 %2d | %s (ID: %s)\n", s.Sequence, s.Name, s.StopID)
	}
	if len(stops) > limit {
		fmt.Printf("  ... 共 %d 站\n", len(stops))
	}

	// Step 3: GetETA
	fmt.Printf("\n=== Step 3: GetETA('%s', direction=0) ===\n", route.RouteID)
	etas, err := p.GetETA(ctx, "", route.RouteID, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetETA 失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("找到 %d 個站的到站資訊\n", len(etas))
	limit = 5
	if len(etas) < limit {
		limit = len(etas)
	}
	for i := 0; i < limit; i++ {
		e := etas[i]
		busInfo := ""
		if len(e.Buses) > 0 {
			busInfo = " 車牌: " + e.Buses[0].PlateNumb
		}
		fmt.Printf("  站序 %2d | ETA: %4d秒 (%2d分)%s\n", e.Sequence, e.ETA, e.ETA/60, busInfo)
	}

	// JSON summary
	fmt.Println("\n=== JSON 摘要 ===")
	summary := map[string]any{
		"keyword":    keyword,
		"routes":     len(routes),
		"firstRoute": route,
		"stops":      len(stops),
		"etas":       len(etas),
	}
	data, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(data))

	fmt.Println("\n✓ SearchRoutes → GetStops → GetETA 完整流程通過")
}
