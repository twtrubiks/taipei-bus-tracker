## Why

CLI 快捷設定（notify.yaml）和 Web 收藏（localStorage）都存了 provider-specific 的 route_id / stop_id。切換 `--provider` 後，舊的 ID 在新 provider 查無資料，收藏/快捷形同作廢。經實測兩個 provider 的站名完全一致，可用名稱反查解決跨 provider ID 映射。

## What Changes

- **BREAKING** CLI 快捷設定格式變更：`route_id` / `stop_id` 改為 `tdx_route_id` / `ebus_route_id` 雙 ID 結構，需遷移舊 notify.yaml
- Web 收藏的 Favorite 型別同步擴充為雙 ID 結構
- 新增 lazy resolve 機制：建立收藏時只存當前 provider 的 ID，首次切換 provider 時用 routeName + stopName 反查並補齊第二組 ID，之後切換不再查詢
- auto 模式 fallback 時順便補齊第二組 ID（免費取得）
- SearchRoutes / GetStops 回傳結果需標記資料來源（Source），以便判斷要存到哪組 ID

## Capabilities

### New Capabilities
- `lazy-dual-id`: 跨 provider 雙 ID 儲存與 lazy resolve 機制，涵蓋 CLI 快捷與 Web 收藏

### Modified Capabilities
- `notify-shortcut`: 快捷設定資料結構從單 ID 改為雙 ID，含舊格式遷移

## Impact

- `cmd/notify/shortcut.go`：Shortcut struct 改為雙 ID + 遷移邏輯
- `web/src/api/types.ts`：Favorite 型別擴充
- `web/src/hooks/useFavorites.ts`：雙 ID 儲存 + localStorage 遷移
- `internal/model/model.go`：Route / Stop 可能需加 Source 欄位
- `internal/handler/fallback.go`：fallback 時標記來源以便補齊 ID
- `notify.yaml`：格式變更（需遷移）
