## 第一段：Go 後端 + CLI

### 1. Model 層 Source 標記

- [x] 1.1 `internal/model/model.go`：Route struct 加 `Source string` 欄位
- [x] 1.2 `internal/model/model.go`：Stop struct 加 `Source string` 欄位
- [x] 1.3 `internal/tdx/`：convertRoutes / convertStops 回傳時填入 Source="tdx"
- [x] 1.4 `internal/ebus/`：相關轉換函式回傳時填入 Source="ebus"

### 2. CLI 快捷雙 ID 結構

- [x] 2.1 `cmd/notify/shortcut.go`：Shortcut struct 改為雙 ID（TDXRouteID / TDXStopID / EBusRouteID / EBusStopID），保留 routeName / stopName
- [x] ~~2.2 `cmd/notify/shortcut.go`：新增遷移函式~~ — 已移除，不需要舊格式遷移
- [x] 2.3 `cmd/notify/shortcut.go`：promptSaveShortcut 依當前 provider source 存入正確的 ID 欄位

### 3. CLI Lazy Resolve

- [x] 3.1 `cmd/notify/shortcut.go`：新增 resolveShortcutID 函式，用 routeName 搜尋 + stopName 比對，取得當前 provider 的 ID
- [x] 3.2 `cmd/notify/main.go`：載入快捷後檢查當前 provider 的 ID 是否為空，為空則呼叫 resolveShortcutID 並回寫 notify.yaml
- [x] 3.3 反查失敗時輸出警告，不阻斷（保留已有 ID 可用）

### 4. Auto Fallback 補齊 ID（Go 端）

- [x] 4.1 `internal/handler/fallback.go`：fallback 成功時在回傳結果標記 Source
- [x] 4.2 CLI 端：fallback 回應時，若第二組 ID 為空，順便存入

### 5. 單元測試（Go）

- [x] 5.1 `internal/provider/provider_test.go`：Build() 三種模式的正常路徑測試（auto+有key、auto+無key、tdx、ebus）
- [x] 5.2 `internal/provider/provider_test.go`：Build() 錯誤路徑測試（tdx 缺憑證、無效 mode）
- [x] ~~5.3-5.5 遷移函式測試~~ — 已移除，不需要舊格式遷移
- [x] 5.6 `cmd/notify/shortcut_test.go`：lazy resolve 測試 — 用 mockProvider 驗證名稱反查流程（搜尋路線 → 比對站名 → 取得 ID）
- [x] 5.7 `cmd/notify/shortcut_test.go`：lazy resolve 測試 — 反查失敗時不阻斷，保留已有 ID
- [x] 5.8 `internal/handler/fallback_test.go`：fallback 成功時回傳結果包含正確 Source 標記

### 6. 整合冒煙測試（go test -tags=integration）

- [x] 6.1 `internal/ebus/ebus_integration_test.go`：實際打 eBus SearchRoutes API，驗證回傳結果不為空且 regex 解析正確
- [x] 6.2 `internal/ebus/ebus_integration_test.go`：實際打 eBus GetStops API，驗證站點列表不為空且欄位完整
- [x] 6.3 `internal/ebus/ebus_integration_test.go`：實際打 eBus GetETA API，驗證回傳 JSON 可解析
- [x] 6.4 `internal/tdx/tdx_integration_test.go`：實際打 TDX SearchRoutes / GetStops / GetETA，驗證格式未變
- [x] 6.5 整合測試加 build tag `//go:build integration`，確保 `go test ./...` 不會觸發，需 `go test -tags=integration ./...` 才跑

### 7. 第一段驗證

- [x] 7.1 `go build ./...` 編譯通過
- [x] 7.2 `go test ./...` 全部通過
- [x] 7.3 `go test -tags=integration ./...` 整合測試通過
- [x] 7.4 手動測試：用 TDX 建快捷 → 切 ebus → 載入快捷 → 確認 lazy resolve 成功並回寫
- [x] 7.5 手動測試：再切回 auto → 確認直接使用已存 ID，不再查詢
- [x] ~~7.6 手動測試：舊 notify.yaml 自動遷移~~ — 不需要，刪除 notify.yaml 重建即可

---

## 第二段：Web 前端

### 8. Web 收藏雙 ID 結構

- [x] 8.1 `web/src/api/types.ts`：Favorite 型別擴充為雙 ID（tdxRouteId / ebusRouteId / tdxStopId / ebusStopId）
- [x] 8.2 ~~`web/src/hooks/useFavorites.ts`：localStorage 讀取時偵測舊格式，自動遷移為雙 ID~~ — 已簡化，不需遷移（清除 localStorage 重新收藏即可）
- [x] 8.3 `web/src/hooks/useFavorites.ts`：addFavorite 依 source 存入正確的 ID 欄位
- [x] 8.4 `web/src/hooks/useFavoritesEta.ts`：查詢 ETA 時依當前 provider 選用正確的 ID
- [x] 8.5 Web 端：fallback 回應時，若第二組 ID 為空，順便存入

### 9. 單元測試（Web）

- [x] ~~9.1 `web/src/hooks/useFavorites.test.ts`：localStorage 舊格式遷移測試~~ — 已移除，不需要舊格式遷移（與 CLI 端策略一致）

### 10. 第二段驗證

- [x] 10.1 `go build ./...` 編譯通過
- [x] 10.2 `go test ./...` 全部通過
- [x] 10.3 Web `npm run lint` + `npm test` 通過
- [x] 10.4 手動測試：Web 收藏在切換 provider 後仍能正常顯示 ETA
