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

	// 印出前 10 個站
	limit := 10
	if len(etas) < limit {
		limit = len(etas)
	}
	for i := 0; i < limit; i++ {
		e := etas[i]
		busInfo := "無車輛"
		if len(e.Buses) > 0 {
			plates := ""
			for j, b := range e.Buses {
				if j > 0 {
					plates += ", "
				}
				plates += b.PlateNumb
			}
			busInfo = "車牌: " + plates
		}
		fmt.Printf("  站序 %2d | ETA: %4d 秒 (%2d 分) | %s\n", e.Sequence, e.ETA, e.ETA/60, busInfo)
	}

	fmt.Printf("\n--- 完整 JSON (前 3 筆) ---\n")
	limit = 3
	if len(etas) < limit {
		limit = len(etas)
	}
	data, _ := json.MarshalIndent(etas[:limit], "", "  ")
	fmt.Println(string(data))
}
