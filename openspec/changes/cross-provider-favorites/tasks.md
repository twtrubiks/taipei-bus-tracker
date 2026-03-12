## 第一段：Go 後端 + CLI

### 1. Model 層 Source 標記

- [ ] 1.1 `internal/model/model.go`：Route struct 加 `Source string` 欄位
- [ ] 1.2 `internal/model/model.go`：Stop struct 加 `Source string` 欄位
- [ ] 1.3 `internal/tdx/`：convertRoutes / convertStops 回傳時填入 Source="tdx"
- [ ] 1.4 `internal/ebus/`：相關轉換函式回傳時填入 Source="ebus"

### 2. CLI 快捷雙 ID 結構

- [ ] 2.1 `cmd/notify/shortcut.go`：Shortcut struct 改為雙 ID（TDXRouteID / TDXStopID / EBusRouteID / EBusStopID），保留 routeName / stopName
- [ ] 2.2 `cmd/notify/shortcut.go`：新增遷移函式，載入 notify.yaml 時偵測舊格式（單 route_id），根據 ID 前綴判斷來源並轉換為雙 ID，回寫檔案
- [ ] 2.3 `cmd/notify/shortcut.go`：promptSaveShortcut 依當前 provider source 存入正確的 ID 欄位

### 3. CLI Lazy Resolve

- [ ] 3.1 `cmd/notify/shortcut.go`：新增 resolveShortcutID 函式，用 routeName 搜尋 + stopName 比對，取得當前 provider 的 ID
- [ ] 3.2 `cmd/notify/main.go`：載入快捷後檢查當前 provider 的 ID 是否為空，為空則呼叫 resolveShortcutID 並回寫 notify.yaml
- [ ] 3.3 反查失敗時輸出警告，不阻斷（保留已有 ID 可用）

### 4. Auto Fallback 補齊 ID（Go 端）

- [ ] 4.1 `internal/handler/fallback.go`：fallback 成功時在回傳結果標記 Source
- [ ] 4.2 CLI 端：fallback 回應時，若第二組 ID 為空，順便存入

### 5. 單元測試（Go）

- [ ] 5.1 `internal/provider/provider_test.go`：Build() 三種模式的正常路徑測試（auto+有key、auto+無key、tdx、ebus）
- [ ] 5.2 `internal/provider/provider_test.go`：Build() 錯誤路徑測試（tdx 缺憑證、無效 mode）
- [ ] 5.3 `cmd/notify/shortcut_test.go`：遷移函式測試 — TDX 格式舊快捷（TPE 開頭）正確遷移到 tdx 欄位
- [ ] 5.4 `cmd/notify/shortcut_test.go`：遷移函式測試 — eBus 格式舊快捷（純數字）正確遷移到 ebus 欄位
- [ ] 5.5 `cmd/notify/shortcut_test.go`：遷移函式測試 — 已是新格式則不動
- [ ] 5.6 `cmd/notify/shortcut_test.go`：lazy resolve 測試 — 用 mockProvider 驗證名稱反查流程（搜尋路線 → 比對站名 → 取得 ID）
- [ ] 5.7 `cmd/notify/shortcut_test.go`：lazy resolve 測試 — 反查失敗時不阻斷，保留已有 ID
- [ ] 5.8 `internal/handler/fallback_test.go`：fallback 成功時回傳結果包含正確 Source 標記

### 6. 整合冒煙測試（go test -tags=integration）

- [ ] 6.1 `internal/ebus/ebus_integration_test.go`：實際打 eBus SearchRoutes API，驗證回傳結果不為空且 regex 解析正確
- [ ] 6.2 `internal/ebus/ebus_integration_test.go`：實際打 eBus GetStops API，驗證站點列表不為空且欄位完整
- [ ] 6.3 `internal/ebus/ebus_integration_test.go`：實際打 eBus GetETA API，驗證回傳 JSON 可解析
- [ ] 6.4 `internal/tdx/tdx_integration_test.go`：實際打 TDX SearchRoutes / GetStops / GetETA，驗證格式未變
- [ ] 6.5 整合測試加 build tag `//go:build integration`，確保 `go test ./...` 不會觸發，需 `go test -tags=integration ./...` 才跑

### 7. 第一段驗證

- [ ] 7.1 `go build ./...` 編譯通過
- [ ] 7.2 `go test ./...` 全部通過
- [ ] 7.3 `go test -tags=integration ./...` 整合測試通過
- [ ] 7.4 手動測試：用 TDX 建快捷 → 切 ebus → 載入快捷 → 確認 lazy resolve 成功並回寫
- [ ] 7.5 手動測試：再切回 auto → 確認直接使用已存 ID，不再查詢
- [ ] 7.6 手動測試：舊 notify.yaml 自動遷移為雙 ID 格式

---

## 第二段：Web 前端

### 8. Web 收藏雙 ID 結構

- [ ] 8.1 `web/src/api/types.ts`：Favorite 型別擴充為雙 ID（tdxRouteId / ebusRouteId / tdxStopId / ebusStopId）
- [ ] 8.2 `web/src/hooks/useFavorites.ts`：localStorage 讀取時偵測舊格式，自動遷移為雙 ID
- [ ] 8.3 `web/src/hooks/useFavorites.ts`：addFavorite 依 source 存入正確的 ID 欄位
- [ ] 8.4 `web/src/hooks/useFavoritesEta.ts`：查詢 ETA 時依當前 provider 選用正確的 ID
- [ ] 8.5 Web 端：fallback 回應時，若第二組 ID 為空，順便存入

### 9. 單元測試（Web）

- [ ] 9.1 `web/src/hooks/useFavorites.test.ts`：localStorage 舊格式遷移測試

### 10. 第二段驗證

- [ ] 10.1 `go build ./...` 編譯通過
- [ ] 10.2 `go test ./...` 全部通過
- [ ] 10.3 Web `npm run lint` + `npm test` 通過
- [ ] 10.4 手動測試：Web 收藏在切換 provider 後仍能正常顯示 ETA
