package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
	"github.com/twtrubiks/taipei-bus-tracker/internal/tdx"
)

func main() {
	clientID := os.Getenv("TDX_CLIENT_ID")
	clientSecret := os.Getenv("TDX_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		fmt.Fprintln(os.Stderr, "請設定環境變數 TDX_CLIENT_ID 和 TDX_CLIENT_SECRET")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	city := "Taipei"
	keyword := "299"
	direction := 0

	fmt.Println("=== TDX API 整合測試 ===")
	fmt.Printf("City: %s, 關鍵字: %s, 方向: %d (去程)\n\n", city, keyword, direction)

	// Step 0: 先測 OAuth token
	fmt.Println("--- Step 0: 測試 OAuth Token ---")
	token, err := testOAuth(ctx, clientID, clientSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "OAuth 失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Token 取得成功 (前 20 字元: %s...)\n\n", token[:min(20, len(token))])

	// Step 1: 測原始 API 回傳格式（印出 raw JSON）
	fmt.Println("--- Step 1: 原始 JSON 回傳格式 ---")
	testRawAPI(ctx, token, city, keyword, direction)

	// Step 2: 透過 Provider 測試轉換後的結果
	fmt.Println("\n(等待 2 秒避免 rate limit...)")
	time.Sleep(2 * time.Second)

	p := tdx.NewProvider(clientID, clientSecret)

	fmt.Println("\n--- Step 2: SearchRoutes ---")
	routes, err := p.SearchRoutes(ctx, city, keyword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SearchRoutes 錯誤: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("找到 %d 條路線\n", len(routes))
	for i, r := range routes {
		if i >= 5 {
			fmt.Printf("  ... 還有 %d 條\n", len(routes)-5)
			break
		}
		fmt.Printf("  [%d] ID=%s 名稱=%s (%s → %s)\n", i, r.RouteID, r.Name, r.StartStop, r.EndStop)
	}

	if len(routes) == 0 {
		fmt.Println("找不到路線，無法繼續測試")
		os.Exit(0)
	}

	// 用第一個搜尋結果繼續測試
	routeID := routes[0].RouteID
	fmt.Printf("\n使用路線: %s (%s) 繼續測試\n", routes[0].Name, routeID)

	time.Sleep(1 * time.Second)
	fmt.Println("\n--- Step 3: GetStops ---")
	stops, err := p.GetStops(ctx, city, routeID, direction)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetStops 錯誤: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("找到 %d 個站點\n", len(stops))
	limit := min(10, len(stops))
	for i := 0; i < limit; i++ {
		s := stops[i]
		fmt.Printf("  站序 %2d | ID=%s | %s\n", s.Sequence, s.StopID, s.Name)
	}
	if len(stops) > limit {
		fmt.Printf("  ... 還有 %d 個站\n", len(stops)-limit)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("\n--- Step 4: GetETA ---")
	etas, err := p.GetETA(ctx, city, routeID, direction)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetETA 錯誤: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("取得 %d 個站點的到站資訊\n", len(etas))
	limit = min(10, len(etas))
	for i := 0; i < limit; i++ {
		e := etas[i]
		status := model.ETAStatus(e.ETA)
		busInfo := "無車輛"
		if len(e.Buses) > 0 {
			plates := make([]string, len(e.Buses))
			for j, b := range e.Buses {
				plates[j] = b.PlateNumb
			}
			busInfo = "車牌: " + strings.Join(plates, ", ")
		}
		fmt.Printf("  站序 %2d | %-6s | ETA=%5d秒 | %s | %s\n", e.Sequence, e.StopName, e.ETA, status, busInfo)
	}

	// 印出轉換後 JSON
	fmt.Printf("\n--- 轉換後 JSON (前 3 筆 ETA) ---\n")
	jsonLimit := min(3, len(etas))
	data, _ := json.MarshalIndent(etas[:jsonLimit], "", "  ")
	fmt.Println(string(data))

	fmt.Println("\n=== 測試完成 ===")
}

func testOAuth(ctx context.Context, clientID, clientSecret string) (string, error) {
	tokenURL := "https://tdx.transportdata.tw/auth/realms/TDXConnect/protocol/openid-connect/token"
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP 請求失敗: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("JSON 解析失敗: %w (%s)", err, body)
	}

	fmt.Printf("  Token Type: %s\n", result.TokenType)
	fmt.Printf("  Expires In: %d 秒\n", result.ExpiresIn)
	return result.AccessToken, nil
}

func buildRawURL(base, filter string) string {
	params := url.Values{}
	params.Set("$filter", filter)
	params.Set("$format", "JSON")
	return base + "?" + params.Encode()
}

func testRawAPI(ctx context.Context, token, city, keyword string, direction int) {
	baseURL := "https://tdx.transportdata.tw/api/basic/v2/Bus"
	client := &http.Client{Timeout: 10 * time.Second}

	// Raw SearchRoutes
	routeURL := buildRawURL(baseURL+"/Route/City/"+city,
		fmt.Sprintf("contains(RouteName/Zh_tw,'%s')", keyword))
	routeURL += "&" + url.Values{"$top": {"3"}}.Encode()
	fmt.Printf("\n[SearchRoutes] URL: %s\n", routeURL)
	rawBody := fetchRaw(ctx, client, token, routeURL)
	if rawBody != nil {
		printPrettyJSON("SearchRoutes 原始回傳 (前 3 筆)", rawBody)

		var rawRoutes []json.RawMessage
		if json.Unmarshal(rawBody, &rawRoutes) == nil && len(rawRoutes) > 0 {
			var first map[string]interface{}
			if json.Unmarshal(rawRoutes[0], &first) == nil {
				if uid, ok := first["RouteUID"].(string); ok {
					routeID := uid

					time.Sleep(1 * time.Second) // 避免 rate limit

					// Raw GetStops
					stopsURL := buildRawURL(baseURL+"/StopOfRoute/City/"+city,
						fmt.Sprintf("RouteUID eq '%s' and Direction eq %d", routeID, direction))
					stopsURL += "&" + url.Values{"$top": {"1"}}.Encode()
					fmt.Printf("\n[GetStops] URL: %s\n", stopsURL)
					rawStops := fetchRaw(ctx, client, token, stopsURL)
					if rawStops != nil {
						printPrettyJSON("StopOfRoute 原始回傳 (前 1 筆)", rawStops)
					}

					time.Sleep(1 * time.Second)

					// Raw GetETA
					etaURL := buildRawURL(baseURL+"/EstimatedTimeOfArrival/City/"+city,
						fmt.Sprintf("RouteUID eq '%s' and Direction eq %d", routeID, direction))
					etaURL += "&" + url.Values{"$top": {"5"}}.Encode()
					fmt.Printf("\n[GetETA] URL: %s\n", etaURL)
					rawETA := fetchRaw(ctx, client, token, etaURL)
					if rawETA != nil {
						printPrettyJSON("ETA 原始回傳 (前 5 筆)", rawETA)
					}
				}
			}
		}
	}
}

func fetchRaw(ctx context.Context, client *http.Client, token, url string) []byte {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  建立請求失敗: %v\n", err)
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  HTTP 請求失敗: %v\n", err)
		return nil
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "  HTTP %d: %s\n", resp.StatusCode, string(body)[:min(200, len(body))])
		return nil
	}
	return body
}

func printPrettyJSON(label string, data []byte) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("  %s (raw): %s\n", label, string(data)[:min(500, len(data))])
		return
	}
	pretty, _ := json.MarshalIndent(v, "  ", "  ")
	fmt.Printf("  %s:\n  %s\n", label, string(pretty))
}
