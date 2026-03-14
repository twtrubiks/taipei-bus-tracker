## Why

每天查公車到站時間，現有工具不是手機優先就是有廣告。需要一個自架的、無廣告的、電腦手機都好用的台北公車即時到站查詢工具。部署在自己的 Linux server 上，完全自己掌控。

## What Changes

- 新建 Go 後端 API proxy，整合 TDX（主要）與 eBus（備援）雙資料源
- 新建 React PWA 前端，支援電腦瀏覽器 + 手機瀏覽器（可加到主畫面）
- 支援路線搜尋、方向選擇、指定站點即時到站時間查詢
- 每 15 秒自動更新，與官方資料同步
- 收藏常用路線/站點，首頁直接顯示即時狀態
- 深色模式、到站瀏覽器通知提醒

## Capabilities

### New Capabilities

- `go-api-proxy`: Go 後端 API proxy，處理 TDX/eBus 認證、資料統一、快取、fallback 機制
- `bus-data-source`: 雙資料源抽象層，TDX 為主要來源、eBus 為備援，統一資料模型
- `pwa-frontend`: React PWA 前端，包含路線搜尋、到站顯示、收藏、深色模式、通知提醒
- `deployment`: Linux server 部署方案，Go binary + 靜態檔案 serve

### Modified Capabilities

（無既有 capability，全新專案）

## Impact

- **新增依賴**: Go standard library + 少量第三方（router）、React + TypeScript + Vite
- **外部 API**: TDX API（需註冊取得 Client ID/Secret）、ebus.gov.taipei（CSRF token）
- **部署**: Linux server 上運行單一 Go binary，serve API + 靜態前端
- **網路**: server 需能對外連線到 TDX 和 ebus API
