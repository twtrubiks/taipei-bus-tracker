package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/twtrubiks/taipei-bus-tracker/internal/ebus"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	p := ebus.NewProvider()

	// 299 路公車
	routeID := "0100029900"
	direction := 0

	fmt.Println("=== eBus API 整合測試 ===")
	fmt.Printf("路線 ID: %s (299路), 方向: %d (去程)\n\n", routeID, direction)

	etas, err := p.GetETA(ctx, "", routeID, direction)
	if err != nil {
		fmt.Fprintf(os.Stderr, "錯誤: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功取得 %d 個站點的到站資訊\n\n", len(etas))

	// 印出所有有車的站：站上（bi）與站間（bo）分開列
	atStop, between := 0, 0
	for _, e := range etas {
		if len(e.Buses) == 0 && len(e.DepartedBuses) == 0 {
			continue
		}
		fmt.Printf("  站序 %2d | ETA: %4d 秒 (%2d 分) | 站上: %-24s | 站間: %s\n",
			e.Sequence, e.ETA, e.ETA/60, plates(e.Buses), plates(e.DepartedBuses))
		atStop += len(e.Buses)
		between += len(e.DepartedBuses)
	}
	fmt.Printf("\n車輛統計：站上 %d 台、站間 %d 台，合計 %d 台\n", atStop, between, atStop+between)

	fmt.Printf("\n--- 完整 JSON (前 3 筆) ---\n")
	limit := 3
	if len(etas) < limit {
		limit = len(etas)
	}
	data, _ := json.MarshalIndent(etas[:limit], "", "  ")
	fmt.Println(string(data))
}

// plates formats plate numbers for one-line display.
func plates(buses []model.Bus) string {
	if len(buses) == 0 {
		return "-"
	}
	out := make([]string, 0, len(buses))
	for _, b := range buses {
		out = append(out, b.PlateNumb)
	}
	return strings.Join(out, ", ")
}
