## Context

CLI 快捷（notify.yaml）和 Web 收藏（localStorage）都存了 provider-specific 的 route_id / stop_id。TDX 用 `TPE10723` 格式，eBus 用 `0100000100` 格式，切換 provider 後舊 ID 查無資料。

經實測路線 1（20 站）和路線 22（29 站），兩個 provider 的站名完全一致（來自同一套交通部標準），可用名稱做跨 provider 映射。

## Goals / Non-Goals

**Goals:**
- CLI 快捷和 Web 收藏在切換 provider 後仍能正常使用
- 採用 lazy dual-ID 策略：建立時存當前 provider 的 ID，首次切換時用名稱反查補齊第二組 ID，之後切換零查詢
- auto 模式 fallback 時順便補齊第二組 ID（免費取得）
- 舊格式 notify.yaml 和 localStorage 自動遷移

**Non-Goals:**
- 不支援兩個 provider 以外的第三方資料源
- 不處理兩個 provider 站名不一致的 edge case（實測一致）
- 不改 FallbackService 的 fallback 邏輯本身

## Decisions

### 1. 雙 ID 資料結構

**選擇**: 將單一 `route_id` / `stop_id` 拆成 `tdx_route_id` / `tdx_stop_id` + `ebus_route_id` / `ebus_stop_id`。

**替代方案**:
- 只存名稱、每次都反查：太慢，每次啟動都要搜尋 + 匹配
- 存時就撈雙 ID（方案 C）：使用者可能只有一個 provider 的 key，無法撈雙邊

**理由**: lazy 策略兼顧效能和可用性，第一次切換多一次查詢，之後就不用了。

### 2. 反查流程（lazy resolve）

載入收藏/快捷時，如果當前 provider 對應的 ID 欄位為空：
1. 用 `routeName` 呼叫 `SearchRoutes` 找路線
2. 從結果中精確比對 `routeName` 取得新 route_id
3. 用新 route_id 呼叫 `GetStops` 取站點列表
4. 用 `stopName` 比對取得新 stop_id
5. 將新 ID 回寫到儲存（notify.yaml / localStorage）

### 3. Source 標記

**選擇**: 在 `model.Route` 和 `model.Stop` 加 `Source string` 欄位（`"tdx"` / `"ebus"`），與 `StopETA.Source` 一致。FallbackService 在回傳結果時標記來源。

**理由**: auto 模式 fallback 時需要知道回傳的資料來自哪個 provider，才能存到正確的 ID 欄位。`StopETA` 已有 `Source` 欄位可參考。

### 4. 舊格式遷移

**CLI**: 載入 notify.yaml 時偵測舊格式（有 `route_id` 無 `tdx_route_id`），自動轉換並回寫。根據 ID 格式判斷來源：`TPE` 開頭為 TDX，純數字為 eBus。

**Web**: localStorage 讀取時偵測舊格式（有 `routeId` 無 `tdxRouteId`），轉換後存回。同樣根據 ID 格式判斷。

## Risks / Trade-offs

- **站名比對可靠性**: 實測站名完全一致，但未來若任一方改站名可能導致比對失敗 → 比對失敗時保持當前 ID 可用，僅 log 警告，不阻斷使用
- **遷移時 ID 格式判斷**: `TPE` 開頭判定為 TDX 是啟發式規則，理論上可靠但非 100% → 遷移失敗不刪除原資料，使用者可手動修正
- **反查的額外 API 呼叫**: 首次切換需 SearchRoutes + GetStops 兩次呼叫 → 只發生一次，之後 ID 被回寫後不再觸發
