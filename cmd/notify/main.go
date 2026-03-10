package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/cache"
	"github.com/twtrubiks/taipei-bus-tracker/internal/config"
	"github.com/twtrubiks/taipei-bus-tracker/internal/ebus"
	"github.com/twtrubiks/taipei-bus-tracker/internal/handler"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
	"github.com/twtrubiks/taipei-bus-tracker/internal/tdx"
)

const defaultCity = "Taipei"

func main() {
	listFlag := flag.Bool("list", false, "列出所有快捷")
	deleteFlag := flag.String("delete", "", "刪除指定快捷")
	flag.Parse()

	// Handle --list
	if *listFlag {
		listShortcuts()
		return
	}

	// Handle --delete
	if *deleteFlag != "" {
		deleteShortcut(*deleteFlag)
		return
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("設定載入失敗: %v", err)
	}

	var primary model.BusDataSource
	if cfg.TDX.ClientID != "" && cfg.TDX.ClientSecret != "" {
		primary = tdx.NewProvider(cfg.TDX.ClientID, cfg.TDX.ClientSecret)
		fmt.Println("TDX provider 已初始化")
	}

	ebusProvider := ebus.NewProvider()
	var fallbackSrc model.BusDataSource = ebusProvider

	if primary == nil {
		primary = ebusProvider
		fallbackSrc = nil
		fmt.Println("使用 eBus 作為主要資料來源（無 TDX 憑證）")
	}

	c := cache.New(30 * time.Second)
	defer c.Close()

	svc := handler.NewFallbackService(primary, fallbackSrc, c)

	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignal()

	var routeID, routeName, startStop, endStop string
	var direction, threshold int
	var stop model.Stop

	// Check for positional arg (shortcut name)
	if shortcutName := flag.Arg(0); shortcutName != "" {
		s := loadShortcut(shortcutName)
		routeID = s.RouteID
		routeName = s.RouteName
		startStop = s.StartStop
		endStop = s.EndStop
		direction = s.Direction
		stop = model.Stop{StopID: s.StopID, Name: s.StopName, Sequence: s.StopSequence}
		threshold = s.Threshold
	} else {
		scanner := bufio.NewScanner(os.Stdin)

		route, err := selectRoute(ctx, svc, scanner)
		if err != nil {
			log.Fatalf("路線選擇失敗: %v", err)
		}

		direction, err = selectDirection(route, scanner)
		if err != nil {
			log.Fatalf("方向選擇失敗: %v", err)
		}

		stop, err = selectStop(ctx, svc, scanner, route.RouteID, direction)
		if err != nil {
			log.Fatalf("站點選擇失敗: %v", err)
		}

		threshold = selectThreshold(scanner)
		routeID = route.RouteID
		routeName = route.Name
		startStop = route.StartStop
		endStop = route.EndStop

		promptSaveShortcut(scanner, Shortcut{
			RouteID:      routeID,
			RouteName:    routeName,
			StartStop:    startStop,
			EndStop:      endStop,
			Direction:    direction,
			StopID:       stop.StopID,
			StopName:     stop.Name,
			StopSequence: stop.Sequence,
			Threshold:    threshold,
		})
	}

	notifyCmd := detectNotifyTool()
	detectSoundTool()

	dirLabel := fmt.Sprintf("%s：%s→%s", formatDirection(direction), startStop, endStop)
	if direction == 1 {
		dirLabel = fmt.Sprintf("%s：%s→%s", formatDirection(direction), endStop, startStop)
	}

	fmt.Printf("\n✓ 監控中 %s %s（%s），%d 分鐘前通知  Ctrl+C 停止\n",
		routeName, stop.Name, dirLabel, threshold)
	if notifyCmd == "" {
		fmt.Println("⚠ 未偵測到通知工具，僅 terminal 顯示模式")
	}
	fmt.Println("─────────────────────────────────────────────")

	runMonitor(ctx, svc, routeID, direction, stop, routeName, threshold, notifyCmd)
}

// detectNotifyTool returns the notification command to use:
// "notify-send" > "kdialog" > "" (none)
func detectNotifyTool() string {
	if _, err := exec.LookPath("notify-send"); err == nil {
		fmt.Println("通知工具: notify-send")
		return "notify-send"
	}
	if _, err := exec.LookPath("kdialog"); err == nil {
		fmt.Println("通知工具: kdialog (KDE)")
		return "kdialog"
	}
	fmt.Println("⚠ 未偵測到通知工具（notify-send 或 kdialog），桌面通知將無法使用")
	fmt.Println("  安裝方式: sudo apt install libnotify-bin")
	return ""
}

func selectRoute(ctx context.Context, svc *handler.FallbackService, scanner *bufio.Scanner) (model.Route, error) {
	for {
		fmt.Print("搜尋路線: ")
		if !scanner.Scan() {
			return model.Route{}, fmt.Errorf("輸入中斷")
		}
		keyword := strings.TrimSpace(scanner.Text())
		if keyword == "" {
			continue
		}

		routes, err := svc.SearchRoutes(ctx, defaultCity, keyword)
		if err != nil {
			fmt.Printf("搜尋失敗: %v，請重試\n", err)
			continue
		}

		if len(routes) == 0 {
			fmt.Println("查無路線，請重新輸入")
			continue
		}

		if len(routes) == 1 {
			fmt.Printf("自動選擇: %s (%s→%s)\n", routes[0].Name, routes[0].StartStop, routes[0].EndStop)
			return routes[0], nil
		}

		for i, r := range routes {
			fmt.Printf("  [%d] %s (%s→%s)\n", i+1, r.Name, r.StartStop, r.EndStop)
		}

		fmt.Print("選擇: ")
		if !scanner.Scan() {
			return model.Route{}, fmt.Errorf("輸入中斷")
		}
		choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || choice < 1 || choice > len(routes) {
			fmt.Println("無效選擇，請重新搜尋")
			continue
		}
		return routes[choice-1], nil
	}
}

func selectDirection(route model.Route, scanner *bufio.Scanner) (int, error) {
	fmt.Printf("方向:\n  [0] 去程 (%s→%s)\n  [1] 回程 (%s→%s)\n",
		route.StartStop, route.EndStop, route.EndStop, route.StartStop)
	for {
		fmt.Print("選擇 [0]: ")
		if !scanner.Scan() {
			return 0, fmt.Errorf("輸入中斷")
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" || text == "0" {
			return 0, nil
		}
		if text == "1" {
			return 1, nil
		}
		fmt.Println("請輸入 0 或 1")
	}
}

func selectStop(ctx context.Context, svc *handler.FallbackService, scanner *bufio.Scanner, routeID string, direction int) (model.Stop, error) {
	stops, err := svc.GetStops(ctx, defaultCity, routeID, direction)
	if err != nil {
		return model.Stop{}, fmt.Errorf("取得站點失敗: %w", err)
	}
	if len(stops) == 0 {
		return model.Stop{}, fmt.Errorf("無站點資料")
	}

	fmt.Println("站點:")
	for i, s := range stops {
		fmt.Printf("  [%d] %s\n", i+1, s.Name)
	}

	for {
		fmt.Print("選擇: ")
		if !scanner.Scan() {
			return model.Stop{}, fmt.Errorf("輸入中斷")
		}
		choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || choice < 1 || choice > len(stops) {
			fmt.Println("無效選擇，請重新輸入")
			continue
		}
		return stops[choice-1], nil
	}
}

func selectThreshold(scanner *bufio.Scanner) int {
	for {
		fmt.Print("幾分鐘前通知 [5]: ")
		if !scanner.Scan() {
			return 5
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			return 5
		}
		n, err := strconv.Atoi(text)
		if err != nil || n <= 0 {
			fmt.Println("請輸入正整數")
			continue
		}
		return n
	}
}
